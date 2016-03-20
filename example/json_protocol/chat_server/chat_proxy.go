package chat_server

import (
	"github.com/jeidee/gonet"
	. "github.com/jeidee/gonet/example/json_protocol/chat_packet"
)

type ChatProxy struct {
	server *ChatServer
}

func (proxy *ChatProxy) ResLogin(session *gonet.Session, obj ResLogin) {
	session.Send(obj)
}

func (proxy *ChatProxy) ResUserList(session *gonet.Session, obj ResUserList) {
	session.Send(obj)
}

func (proxy *ChatProxy) NotifyChat(session *gonet.Session, obj NotifyChat) {
	session.Send(obj)
}
