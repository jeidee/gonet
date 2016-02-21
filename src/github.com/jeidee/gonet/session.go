package gonet

import (
	"fmt"
	"io"
	"net"
)

type Session struct {
	server   *Server
	conn     net.Conn
	protocol Protocol
}

func NewSession(conn net.Conn, server *Server) *Session {
	session := new(Session)

	session.server = server
	session.protocol = server.protocol.Make(conn)
	session.conn = conn

	session.Info("Accepted connection %s -> %s", conn.RemoteAddr().String(), conn.LocalAddr().String())

	return session
}

/* public functions - APIs */

func (session *Session) Run() {
	go session.recvLoop()
}

func (session *Session) Close() {
	session.server.CloseSession(session)
}

func (session *Session) Send(data interface{}) error {
	_, err := session.protocol.Encode(session, data)
	if err != nil {
		session.Error(err, "Sending data failed.")
	}
	return err
}

func (session *Session) Conn() net.Conn {
	return session.conn
}

func (session *Session) Panic(err error) {
	session.Error(err, "")
	session.Close()
}

func (session *Session) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[DEBUG]", msg)
}

func (session *Session) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[INFO]", msg)
}

func (session *Session) Error(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[ERR]", msg, err)
}

/* private functions */

func (session *Session) recvLoop() {
	for {
		data, err := session.protocol.Decode(session)

		if err != nil {
			if err == io.EOF {
				session.Info("Connection closed by remote host.")
				session.Close()
			} else {
				session.Panic(err)
			}
			return
		}

		session.server.Incomming(session, data)
	}
}
