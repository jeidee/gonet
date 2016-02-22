package main

import (
	"github.com/jeidee/gonet"
)

type EchoSession struct {
	session *gonet.Session

	id int32
}

func NewEchoSession(session *gonet.Session, id int32) *EchoSession {
	echoSession := new(EchoSession)
	echoSession.session = session
	echoSession.id = id

	return echoSession
}

type EchoServer struct {
	gonet.Server

	echoSessions  map[*gonet.Session]*EchoSession
	lastSessionId int32
}

///////////////////////////////////////////////////////////////////////////////
/* Server Event Handlers
... This code belong to the goroutine of server's eventLoop. */

func (s *EchoServer) OnAccept(session *gonet.Session) {
	s.lastSessionId += 1
	s.echoSessions[session] = NewEchoSession(session, s.lastSessionId)
	s.Info("EchoServer ... OnAccept...%v, Cu is %v.", s.lastSessionId, len(s.echoSessions))
}

func (s *EchoServer) OnClose(session *gonet.Session) {
	s.lastSessionId -= 1
	delete(s.echoSessions, session)
	s.Info("EchoServer ... OnClose...%v, Cu is %v.", s.lastSessionId, len(s.echoSessions))
}

func (s *EchoServer) OnIncomming(data *gonet.IncommingData) {
	s.Info("EchoServer ... OnIncomming...%v", data.Data)
	// echo
	err := data.Session.Send(data.Data)
	if err != nil {
		s.Error(err, "Echoing failed.")
	}

}

/* End of Server Event Handlers */
///////////////////////////////////////////////////////////////////////////////

func NewEchoServer(port int16) *EchoServer {
	s := new(EchoServer)

	s.echoSessions = make(map[*gonet.Session]*EchoSession)
	s.lastSessionId = 0
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
