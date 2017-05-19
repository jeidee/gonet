package gonet

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	NetHost

	sessionGroupManager *SessionGroupManager

	lisener   net.Listener
	port      int16
	protocol  Protocol
	isRunning bool

	serverEventHandler ServerEventHandler
	stub               Stub

	accepting chan *Session
	closing   chan *Session
	incomming chan *IncommingData
}

func NewServer(port int16, protocol Protocol, serverEventHandler ServerEventHandler, stub Stub) *Server {
	s := new(Server)

	s.Init(port, protocol, serverEventHandler, stub)

	return s
}

/* Implemetations NetHost interfaceds */

func (self *Server) MakeProtocol(conn net.Conn) Protocol {
	return self.protocol.Make(conn)
}

func (self *Server) Incomming(session *Session, data interface{}) {
	self.incomming <- &IncommingData{session, data}
}

/* public functions - APIs */

func (self *Server) IsRunning() bool {
	return self.isRunning
}

func (self *Server) Init(port int16, protocol Protocol, serverEventHandler ServerEventHandler, stub Stub) {
	self.port = port
	self.protocol = protocol
	self.isRunning = false
	self.serverEventHandler = serverEventHandler
	self.stub = stub

	self.sessionGroupManager = NewSessionGroupManager("root", nil)

	self.accepting = make(chan *Session)
	self.closing = make(chan *Session)
	self.incomming = make(chan *IncommingData)
}

func (self *Server) Run() bool {
	if self.isRunning {
		return false
	}

	go self.eventLoop()

	if !self.acceptLoop() {
		return false
	}

	return true
}

func (self *Server) Close(session *Session) {
	session.conn.Close()
	self.sessionGroupManager.Leave(session)
	self.serverEventHandler.OnClose(session)
	self.Debug("Session is closed!")
}

func (self *Server) Broadcast(data interface{}, exceptSessions ...*Session) {
	self.sessionGroupManager.Broadcast(data, exceptSessions...)
}

func (self *Server) SessionGroupManager() *SessionGroupManager {
	return self.sessionGroupManager
}

func (self *Server) Panic(err error, format string, args ...interface{}) {
	self.raiseError(err, format, args...)
	os.Exit(1)
}

func (self *Server) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[S:DEBUG]", msg)
}

func (self *Server) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[S:INFO]", msg)
}

func (self *Server) Error(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[S:ERR]", msg, err)
}

/* private functions */

func (self *Server) eventLoop() {
	for {
		select {
		case session := <-self.accepting:
			self.sessionGroupManager.Join(session)

			self.serverEventHandler.OnAccept(session)
			self.Debug("New session is accepted!")

		case session := <-self.closing:
			self.Close(session)

		case data := <-self.incomming:
			if data.Data == nil {
				self.Close(data.Session)
				continue
			}

			if self.stub != nil {
				if err := self.stub.ParseAndDo(data); err != nil {
					self.Error(err, "Parsing failed.")
				}
			}

			self.serverEventHandler.OnIncomming(data)
			self.Debug("Incomming new data %v, %v", data.Session, data.Data)
		}
	}

}

func (self *Server) acceptLoop() bool {

	listenPort := fmt.Sprintf(":%d", self.port)
	var err error

	self.lisener, err = net.Listen("tcp", listenPort)
	if err != nil {
		self.Panic(err, "Listening failed!")
		return false
	}
	defer self.lisener.Close()

	self.isRunning = true
	self.Debug("Reusing listening port for %d", self.port)

	for {
		if self.isRunning == false {
			break
		}

		conn, err := self.lisener.Accept()
		if err != nil {
			self.raiseError(err, "Accept failed!")
			continue
		}

		session := NewSession(conn, self)
		self.accepting <- session
		session.Run()
	}

	return true
}

func (self *Server) raiseError(err error, format string, args ...interface{}) {
	self.Error(err, format, args...)
}
