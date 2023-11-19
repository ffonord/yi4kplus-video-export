package ftp

import (
	"github.com/jlaffaye/ftp"
	"net"
	"time"
)

type FTPConnFactory struct {
}

func (c *FTPConnFactory) NewConn(host, port string, timeout time.Duration) (Conn, error) {
	addr := net.JoinHostPort(host, port)
	return ftp.Dial(addr, ftp.DialWithTimeout(timeout))
}
