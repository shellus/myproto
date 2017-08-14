package main

import (
	"net"
	"crypto/rand"
	"log"
	"github.com/shellus/myproto"
)


func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	proto := myproto.NewMyProto(conn)

	randByte := make([]byte, 1024*50)
	randByteN, err := rand.Read(randByte)
	if err != nil {
		panic(err)
	}
	log.Println(randByteN)
	err = proto.Send("hello", myproto.HelloProto{Name: "shellus", Age: 18, RandByte: randByte})
	if err != nil {
		panic(err)
	}
}
