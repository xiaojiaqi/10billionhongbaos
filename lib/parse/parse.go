package parse

import (
	"encoding/binary"
	. "github.com/10billionhongbaos/lib/message"
	"net"
)

func ReadByte(conn net.Conn, buff []byte, length int) bool {
	readed := 0

	for {
		n, err := conn.Read(buff[readed:length])
		if err != nil {
			return false
		}
		readed += n
		if readed == length {
			return true
		}
	}
}
func ReadMessage(conn *net.Conn, ppool *MessagePool) (p *Message, result bool) {

	p = ppool.Pop()
	buff := p.Buff
	result = ReadByte(*conn, buff, 4)
	if result != true {
		return nil, result
	}

	length := binary.BigEndian.Uint16(buff[0:])
	length -= 4
	if (length > 40) || (length < 4) {
		return nil, result
	}
	result = ReadByte(*conn, buff[4:], int(length))
	if p.Unserialization(buff) {
		return p, result
	}

	return nil, result
}

func WriteMessage(conn *net.Conn, p *Message) error {
	b := p.Serialization()
	_, err := (*conn).Write(b[0:p.Bufflen])
	p.Rollback()
	return err

}
