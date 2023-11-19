package telnet

import "bufio"

type BufioReaderFactory struct {
}

func (*BufioReaderFactory) NewReader(conn Conn) Reader {
	return bufio.NewReader(conn)
}
