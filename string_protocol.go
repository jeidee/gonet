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

func (self *StringProtocol) Make(conn net.Conn) Protocol {

	newProtocol := new(StringProtocol)

	newProtocol.reader = bufio.NewReader(conn)
	newProtocol.writer = bufio.NewWriter(conn)

	return newProtocol
}

func (self *StringProtocol) Encode(session *Session, data interface{}) (interface{}, error) {
	str, ok := data.(string)
	if ok {
		_, err := self.writer.WriteString(str)
		if err != nil {
			session.Error(err, "bufio.WriteString(%s) failed.", str)
			return data, err
		}

		err = self.writer.Flush()
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

func (self *StringProtocol) Decode(session *Session) (interface{}, error) {
	line, err := self.reader.ReadString('\n')
	return line, err
}
