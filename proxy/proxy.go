package proxy

import (
	"io"
	"net"
)

type Server interface {
	Name() string
	Addr() string
	Handshake(conn net.Conn) (io.ReadWriter, error)
}
