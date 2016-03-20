package main

import (
	"github.com/jeidee/gonet"
)

///////////////////////////////////////////////////////////////////////////////
/* EchoServer */

type EchoServer struct {
	gonet.Server
}

func NewEchoServer(port int16) *EchoServer {
	s := new(EchoServer)

	s.Init(port, new(gonet.StringProtocol), s, nil)

	return s
}

///////////////////////////////////////////////////////////////////////////////
/* Server Event Handlers
... This code belong to the goroutine of server's eventLoop. */

func (s *EchoServer) OnAccept(session *gonet.Session) {
	s.Info("EchoServer ... OnAccept")
}

func (s *EchoServer) OnClose(session *gonet.Session) {
	s.Info("EchoServer ... OnClose")
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

/* EchoServer */
///////////////////////////////////////////////////////////////////////////////

func main() {

	server := NewEchoServer(12345)
	if !server.Run() {
		server.Info("Starting server failed!")
	}

	server.Info("Bye!")
}
