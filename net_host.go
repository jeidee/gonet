package gonet

import (
	"net"
)

type NetHost interface {
	MakeProtocol(net.Conn) Protocol
	Incomming(*Session, interface{})
}
