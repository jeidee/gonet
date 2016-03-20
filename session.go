package gonet

import (
	"fmt"
	"io"
	"net"
)

type Session struct {
	host     NetHost
	conn     net.Conn
	protocol Protocol
}

func NewSession(conn net.Conn, host NetHost) *Session {
	session := new(Session)

	session.host = host
	session.protocol = host.MakeProtocol(conn)
	session.conn = conn

	session.Debug("New connection %s -> %s", conn.RemoteAddr().String(), conn.LocalAddr().String())

	return session
}

/* public functions - APIs */

func (self *Session) Protocol() interface{} {
	return self.protocol
}

func (self *Session) Run() {
	go self.recvLoop()
}

func (self *Session) Send(data interface{}) error {
	_, err := self.protocol.Encode(self, data)
	if err != nil {
		self.Error(err, "Sending data failed.")
	}
	return err
}

func (self *Session) Conn() net.Conn {
	return self.conn
}

func (self *Session) Panic(err error) {
	self.Error(err, "")
	self.conn.Close()
}

func (self *Session) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[SESSION:DEBUG]", msg)
}

func (self *Session) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[SESSION:INFO]", msg)
}

func (self *Session) Error(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[SESSION:ERR]", msg, err)
}

/* private functions */

func (self *Session) recvLoop() {
	for {
		data, err := self.protocol.Decode(self)

		if err != nil {
			if err == io.EOF {
				self.Debug("Connection closed by remote host.")
				self.conn.Close()
			} else {
				self.Panic(err)
			}
			self.host.Incomming(self, nil)
			return
		}

		self.host.Incomming(self, data)
	}
}
