package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

			go func() {
				defer conn.Close()
				sock5Handshake(conn)
			}()
		}
	}()

	{
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("进程退出")
	}
}

func sock5Handshake(conn net.Conn) error {
	print("sock5 handshake")

	reader := bufio.NewReaderSize(conn, 1)

	// 获取建立连接请求
	var connectRequest [2]byte
	_, err := reader.Read(connectRequest[:])
	if err != nil {
		return fmt.Errorf("sock5连接请求获取失败: %w", err)
	}

	return nil
}
