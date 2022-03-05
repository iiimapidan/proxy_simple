package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	ATypeIPv4   = 0x01
	ATypeDomain = 0x03
	ATypeIPv6   = 0x04
)

type Addr struct {
	Name string // domain
	IP   net.IP
	Port int
}

func (a *Addr) String() string {
	port := strconv.Itoa(a.Port)
	if a.IP == nil {
		return net.JoinHostPort(a.Name, port)
	}
	return net.JoinHostPort(a.IP.String(), port)
}

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
				addr, err := sock5Handshake(conn)
				if err != nil {
					log.Printf("sock5握手失败")
					return
				}

				if addr.IP != nil {
					log.Printf("IP:%v Port:%v", addr.IP.String(), addr.Port)
				} else {
					log.Printf("Domain:%v Port:%v", addr.Name, addr.Port)
				}

				// 直连
				remoteConn, err := net.Dial("tcp", addr.String())
				if err != nil {
					log.Printf("连接 %v失败:%v", remoteConn, err)
					return
				}
				defer remoteConn.Close()

				// 流量转发
				go io.Copy(remoteConn, conn)
				io.Copy(conn, remoteConn)
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

//字节数(大端)组转成int(有符号)
func bytesToIntS(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

func sock5Handshake(conn net.Conn) (*Addr, error) {
	reader := bufio.NewReader(conn)

	// 连接请求
	var connectReq [1 + 1 + 255]byte
	_, err := reader.Read(connectReq[:])
	if err != nil {
		return nil, fmt.Errorf("网络读写失败: %w", err)
	}

	ver := connectReq[0]
	if ver != 0x05 {
		return nil, fmt.Errorf("sock5版本不支持: %w", err)
	}

	// 回复请求
	var connectResp [1 + 1]byte
	connectResp[0] = 0x05
	connectResp[1] = 0x00
	_, err = conn.Write(connectResp[:])
	if err != nil {
		return nil, fmt.Errorf("网络读写失败: %w", err)
	}

	// cmd请求
	var cmdReq [1 + 1 + 1 + 1 + 1]byte
	_, err = reader.Read(cmdReq[:])
	if err != nil {
		return nil, fmt.Errorf("网络读写失败: %w", err)
	}

	cmd := cmdReq[1]
	if cmd != 0x01 {
		return nil, fmt.Errorf("cmd请求无效: %w", err)
	}

	address := &Addr{}

	addrLen := 0
	atype := cmdReq[3]
	if atype == ATypeIPv4 {
		addrLen = net.IPv4len
		address.IP = make(net.IP, addrLen)
	} else if atype == ATypeIPv6 {
		addrLen = net.IPv6len
		address.IP = make(net.IP, addrLen)
	} else if atype == ATypeDomain {
		addrLen = int(cmdReq[4])
	}

	if address.IP != nil {
		_, err = reader.Read(address.IP)
		if err != nil {
			return nil, fmt.Errorf("网络读写失败: %w", err)
		}
	} else {
		tmpAdd := make([]byte, addrLen)
		_, err = reader.Read(tmpAdd)
		if err != nil {
			return nil, fmt.Errorf("网络读写失败: %w", err)
		}

		address.Name = string(tmpAdd)
	}

	var port [2]byte
	_, err = reader.Read(port[:])
	if err != nil {
		return nil, fmt.Errorf("网络读写失败: %w", err)
	}

	address.Port, err = bytesToIntS(port[:])
	if err != nil {
		return nil, fmt.Errorf("获取端口失败: %w", err)
	}

	// 回复cmd请求
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return nil, fmt.Errorf("网络读写失败: %w", err)
	}

	return address, nil
}
