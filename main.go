package main

import (
	"log"
	"net"
)

func main() {
	// 启动sock5服务端
	addr := "127.0.0.1:1081"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("无法监听 %v: %v", addr, err)
	}

	log.Printf("sock5:%v服务监听成功", addr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
		}
	}()
}
