package gonet

import (
	"fmt"
	"net"
	"os"
)

type Client struct {
	NetHost

	session   *Session
	protocol  Protocol
	isRunning bool

	clientEventHandler ClientEventHandler
	stub               Stub

	connecting chan *Session
	closing    chan *Session
	incomming  chan *IncommingData
}

func NewClient(protocol Protocol, clientEventHandler ClientEventHandler, stub Stub) *Client {
	self := new(Client)

	self.Init(protocol, clientEventHandler, stub)

	return self
}

/* Implemetations NetHost interfaceds */

func (self *Client) MakeProtocol(conn net.Conn) Protocol {
	return self.protocol.Make(conn)
}

func (self *Client) Incomming(session *Session, data interface{}) {
	self.incomming <- &IncommingData{session, data}
}

/* public functions - APIs */

func (self *Client) Init(protocol Protocol, clientEventHandler ClientEventHandler, stub Stub) {
	self.protocol = protocol
	self.clientEventHandler = clientEventHandler
	self.stub = stub

	self.connecting = make(chan *Session)
	self.closing = make(chan *Session)
	self.incomming = make(chan *IncommingData)
}

func (self *Client) Connect(ip string, port int16) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	go self.eventLoop()

	self.session = NewSession(conn, self)
	self.session.Run()

	self.connecting <- self.session

	return nil
}

func (self *Client) Close() {
	self.session.conn.Close()
	self.clientEventHandler.OnClose(self.session)
	self.Debug("Session is closed!")
}

func (self *Client) Send(data interface{}) error {
	return self.session.Send(data)
}

func (self *Client) Panic(err error, format string, args ...interface{}) {
	self.raiseError(err, format, args)
	os.Exit(1)
}

func (self *Client) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[C:DEBUG]", msg)
}

func (self *Client) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[C:INFO]", msg)
}

func (self *Client) Error(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[C:ERR]", msg, err)
}

/* private functions */

func (self *Client) eventLoop() {
	for {
		select {
		case session := <-self.connecting:
			self.clientEventHandler.OnConnect(session)
			self.Debug("New session is connected!")

		case <-self.closing:
			self.Close()

		case data := <-self.incomming:
			if data == nil {
				self.Close()
				continue
			}

			if self.stub != nil {
				if err := self.stub.ParseAndDo(data); err != nil {
					self.Error(err, "Parsing failed.")
				}
			}

			self.clientEventHandler.OnIncomming(data)
			self.Debug("Incomming new data %v, %v", data.Session, data.Data)
		}
	}

}

func (self *Client) raiseError(err error, format string, args ...interface{}) {
	self.Error(err, format, args...)
}
