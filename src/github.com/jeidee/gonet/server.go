package gonet

import (
	"fmt"
	"net"
	"os"
)

type IncommingData struct {
	session *Session
	data    interface{}
}

type Server struct {
	lisener   net.Listener
	sessions  map[net.Conn]*Session
	port      int16
	protocol  Protocol
	isRunning bool

	accepting chan *Session
	closing   chan *Session
	incomming chan *IncommingData
}

func NewServer(port int16, protocol Protocol) *Server {
	s := new(Server)

	s.port = port
	s.protocol = protocol
	s.isRunning = false

	s.accepting = make(chan *Session)
	s.closing = make(chan *Session)
	s.incomming = make(chan *IncommingData)

	return s
}

/* public functions - APIs */

func (s *Server) Run() bool {
	if s.isRunning {
		return false
	}

	go s.eventLoop()

	if !s.acceptLoop() {
		return false
	}

	return true
}

func (s *Server) AddSession(session *Session) {
	s.accepting <- session
}

func (s *Server) CloseSession(session *Session) {
	s.closing <- session
}

func (s *Server) Incomming(session *Session, data interface{}) {
	s.incomming <- &IncommingData{session, data}
}

func (s *Server) Panic(err error, format string, args ...interface{}) {
	s.raiseError(err, format, args)
	os.Exit(1)
}

func (s *Server) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[DEBUG]", msg)
}

func (s *Server) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[INFO]", msg)
}

func (s *Server) Error(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[ERR]", msg, err)
}

/* private functions */

func (s *Server) eventLoop() {
	for {
		select {
		case session := <-s.accepting:
			s.sessions[session.conn] = session
			s.Info("New session is accepted!")

		case session := <-s.closing:
			delete(s.sessions, session.conn)
			s.Info("Session is closed!")

		case data := <-s.incomming:
			s.Info("Incomming new data %v, %v", data.session, data.data)
		}
	}

}

func (s *Server) acceptLoop() bool {

	s.sessions = make(map[net.Conn]*Session)

	listenPort := fmt.Sprintf(":%d", s.port)
	var err error

	s.lisener, err = net.Listen("tcp", listenPort)
	if err != nil {
		s.Panic(err, "Listening failed!")
		return false
	}
	defer s.lisener.Close()

	s.isRunning = true
	s.Info("Reusing listening port for %d", s.port)

	for {
		if s.isRunning == false {
			break
		}

		conn, err := s.lisener.Accept()
		if err != nil {
			s.raiseError(err, "Accept failed!")
			continue
		}

		session := NewSession(conn, s)
		s.AddSession(session)
		session.Run()
	}

	return true
}

func (s *Server) raiseError(err error, format string, args ...interface{}) {
	s.Error(err, format, args...)
}
