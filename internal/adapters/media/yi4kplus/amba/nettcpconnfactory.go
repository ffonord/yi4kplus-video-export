package amba

import (
	"net"
)

type NetTCPConnFactory struct {
}

func (*NetTCPConnFactory) NewConn(host, port string) (Conn, error) {
	addr := net.JoinHostPort(host, port)
	return net.Dial("tcp", addr)
}
