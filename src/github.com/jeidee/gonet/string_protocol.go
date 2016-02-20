package gonet

import (
	"bufio"
	"errors"
	"net"
)

type StringProtocol struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func (p *StringProtocol) Make(conn net.Conn) Protocol {

	newProtocol := new(StringProtocol)

	newProtocol.reader = bufio.NewReader(conn)
	newProtocol.writer = bufio.NewWriter(conn)

	return newProtocol
}

func (p *StringProtocol) Encode(session *Session, data interface{}) (interface{}, error) {
	str, ok := data.(string)
	if ok {
		_, err := p.writer.WriteString(str)
		if err != nil {
			session.Error(err, "bufio.WriteString(%s) failed.", str)
			return data, err
		}

		err = p.writer.Flush()
		if err != nil {
			session.Error(err, "bufio.Flush() failed.")
			return data, err
		}
	} else {
		err := errors.New("data is not a string!")
		session.Error(err, "")
		return data, err
	}

	return data, nil
}

func (p *StringProtocol) Decode(session *Session) (interface{}, error) {
	line, err := p.reader.ReadString('\n')
	return line, err
}
