package main

import (
	"awesomeRpc/codec"
	"awesomeRpc/server"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on ", l.Addr())
	addr <- l.Addr().String()
	server.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	_ = json.NewEncoder(conn).Encode(server.DefaultOption)
	cc := codec.NewGobCodec(conn)
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("rpc req %d", i))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply: ", reply)
	}
}
