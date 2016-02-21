package main

import (
	"github.com/jeidee/gonet"
)

type EchoSession struct {
	session *gonet.Session

	id int32
}

func NewEchoSession(session *gonet.Session) *EchoSession {
	echoSession := new(EchoSession)
	echoSession.session = session

	return echoSession
}

type EchoServer struct {
	gonet.Server

	echoSessions map[*gonet.Session]*EchoSession
}

func (s *EchoServer) OnAccept(session *gonet.Session) {
	s.echoSessions[session] = NewEchoSession(session)
	s.Info("EchoServer ... OnAccept...Cu is %v.", len(s.echoSessions))
}

func (s *EchoServer) OnClose(session *gonet.Session) {
	delete(s.echoSessions, session)
	s.Info("EchoServer ... OnClose...Cu is %v.", len(s.echoSessions))
}

func (s *EchoServer) OnIncomming(data *gonet.IncommingData) {
	s.Info("EchoServer ... OnIncomming...%v", data.Data)
	// echo
	err := data.Session.Send(data.Data)
	if err != nil {
		s.Error(err, "Echoing failed.")
	}

}

func NewEchoServer(port int16) *EchoServer {
	s := new(EchoServer)

	s.echoSessions = make(map[*gonet.Session]*EchoSession)
	s.Init(port, new(gonet.StringProtocol), s)

	return s
}

func main() {

	server := NewEchoServer(12345)
	if !server.Run() {
		server.Info("Starting server failed!")
	}

	server.Info("Bye!")
}
