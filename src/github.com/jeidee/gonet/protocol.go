package gonet

import (
	"net"
)

type Protocol interface {
	Make(conn net.Conn) Protocol
	Encode(session *Session, data interface{}) (interface{}, error)
	Decode(session *Session) (interface{}, error)
}
