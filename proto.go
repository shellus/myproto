package myproto

import (
	"os"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
	"log"
	"errors"
	"io"
)
func init() {
	gob.Register(HelloProto{})
}
type HelloProto struct {
	Name     string
	Age      int
	RandByte []byte
}
type MyProto struct {
	conn   net.Conn
	logger *log.Logger
}

func NewMyProto(conn net.Conn) *MyProto {
	proto := &MyProto{
		conn:   conn,
		logger: log.New(os.Stderr, "", log.LstdFlags|log.Llongfile),
	}
	return proto
}
func (t *MyProto) Read() error {

	for {
		headerBytes := make([]byte, 128)
		// 读出包长
		headerN, err := io.ReadAtLeast(t.conn, headerBytes, 8)
		if err != nil {
			return errors.New(fmt.Sprintf("Receive header err: %s", err))
		}
		t.logger.Printf("Receive headerN %d", headerN)

		var packageLen uint64 = 0
		err = binary.Read(bytes.NewReader(headerBytes[:8]), binary.BigEndian, &packageLen)
		if err != nil {
			return errors.New(fmt.Sprintf("Read packageLen err: %s", err))
		}
		t.logger.Printf("Parsed the Message len %d", packageLen)

		packageBuffer := bytes.NewBuffer(headerBytes[8:headerN])
		packageBytes := make([]byte, packageLen)
		// 读取包长（注意减去首次接收到的余量）
		packageN, err := io.ReadAtLeast(t.conn, packageBytes, int(packageLen)-(headerN-8))
		if err != nil {
			return errors.New(fmt.Sprintf("Receive packageN err: %s", err))
		}
		packageBuffer.Write(packageBytes[:packageN])
		t.logger.Printf("Received Message len %d", packageBuffer.Len())
	}

	// fix
	return nil
}
func (t *MyProto) Send(messageName string, data interface{}) error {
	tmp := bytes.NewBuffer([]byte{})

	// 写入类型
	err := binary.Write(tmp, binary.BigEndian, uint8(len(messageName)))
	if err != nil {
		return err
	}
	_, err = tmp.WriteString(messageName)
	if err != nil {
		panic(err)
	}

	enc := gob.NewEncoder(tmp)
	err = enc.Encode(data)
	if err != nil {
		return err
	}

	err = binary.Write(t.conn, binary.BigEndian, uint64(tmp.Len()))
	t.logger.Printf("package len %d", tmp.Len())
	if err != nil {
		return err
	}
	n, err := t.conn.Write(tmp.Bytes())
	if err != nil {
		return err
	}
	t.logger.Printf("Write package len %d", n)
	if n != tmp.Len() {
		return errors.New(fmt.Sprintf("Write len Not long enough, total %d write %d", tmp.Len(), n))
	}
	return nil
}
