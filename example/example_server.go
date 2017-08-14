package main

import (
	"net"
	"github.com/shellus/myproto"
	"log"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for{
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn){
			defer conn.Close()

			proto := myproto.NewMyProto(conn)
			//proto.OnMessage("hello", func(message interface{}) {
			//
			//})
			err := proto.Read()
			log.Println(err)
		}(conn)
	}
}