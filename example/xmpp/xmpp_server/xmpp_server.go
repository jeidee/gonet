package xmpp_server

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/jeidee/gonet"
	"github.com/nu7hatch/gouuid"
)

type XmppServer struct {
	gonet.Server

	xmppSessions map[*gonet.Session]*XmppSession
	loginUsers   *gonet.SessionGroupManager
}

func NewXmppServer(port int16) *XmppServer {
	s := new(XmppServer)

	s.xmppSessions = make(map[*gonet.Session]*XmppSession)

	s.Init(port, new(XmppProtocol), s, nil)

	loginUsers, err := s.SessionGroupManager().NewGroup("loginUsers")
	s.Debug("loginUsers %v, err %v", loginUsers, err)
	if err != nil {
		s.Debug("Can't create a new group[%v].", "loginUsers")
		return nil
	}
	s.loginUsers = loginUsers

	return s
}

///////////////////////////////////////////////////////////////////////////////
/* APIs */

/* APIs */
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
/* Server Event Handlers
... This code belong to the goroutine of server's eventLoop. */

func (self *XmppServer) OnAccept(session *gonet.Session) {
	self.Debug("XmppServer ... OnAccept")

	self.xmppSessions[session] = NewXmppSession(session)

}

func (self *XmppServer) OnClose(session *gonet.Session) {
	self.Debug("XmppServer ... OnClose")

	delete(self.xmppSessions, session)
	self.loginUsers.Leave(session)
}

func (self *XmppServer) OnIncomming(data *gonet.IncommingData) {
	self.Debug("XmppServer ... OnIncomming...%v", data.Data)

	xmppSession := self.xmppSessions[data.Session]

	if xmppSession == nil {
		self.Error(errors.New("xmppSession is nil"), "")
		return
	}

	switch t := data.Data.(type) {
	case xml.StartElement:
		self.startStream(xmppSession, t)
	case *saslAuth:
		self.saslAuth(xmppSession, t)
	case *clientIQ:
		self.clientIQ(xmppSession, t)
	case *clientPresence:
		self.clientPresence(xmppSession, t)
	case *clientMessage:
		self.clientMessage(xmppSession, t)
	default:
		self.Error(errors.New("This type is not supported"), "[%v]", t)
	}
}

/* End of Server Event Handlers */
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
/* Internal Functions */
func (self *XmppServer) startStream(xsession *XmppSession, se xml.StartElement) {

	// Parse a start element
	for _, v := range se.Attr {
		switch v.Name.Local {
		case "to":
			xsession.host = v.Value
			break
		case "version":
			xsession.version = v.Value
			break
		}
	}

	self.sendStartStream(xsession)

	switch xsession.state {
	case sessionStateNone:
		//xsession.state = sessionStateTlsNegotiation
		xsession.state = sessionStateSaslNegotiation
		self.sendAuthFeatures(xsession)
		break
	case sessionStateSaslNegotiation:
		xsession.state = sessionStateResourceBinding
		self.sendBindFeatures(xsession)
		break
	}

}

func (self *XmppServer) saslAuth(xsession *XmppSession, sa *saslAuth) {
	self.Debug("saslAuth Mechanism %v %v", sa.Mechanism, sa)
	switch sa.Mechanism {
	case "PLAIN":
		self.responseAuthPlain(xsession, sa)
		break
	}
}

func (self *XmppServer) clientIQ(xsession *XmppSession, iq *clientIQ) {
	self.Debug("clientIQ %v", iq)

	if iq.Bind != (bindBind{}) {
		if iq.Bind.Resource == "" {
			xsession.resource = self.makeResource()
			xsession.state = sessionStateStartSession
		} else {
			xsession.resource = iq.Bind.Resource
			xsession.state = sessionStateStartSession
		}

		self.responseBindingResource(xsession, iq.ID)
		return
	} else if iq.RosterQuery != (rosterQuery{}) {

	}

	self.responseIqUnavailable(xsession, iq)
}

func (self *XmppServer) clientPresence(xsession *XmppSession, presence *clientPresence) {
	self.Debug("clientPresence %v", presence)
}

func (self *XmppServer) clientMessage(xsession *XmppSession, message *clientMessage) {
	self.Debug("clientMessage %v", message)
}

func (self *XmppServer) sendStartStream(xsession *XmppSession) {
	xml := fmt.Sprintf("<?xml version='1.0'?><stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' id='%v' from='%v' version='%v' xml:lang='en'>",
		xsession.id, xsession.host, xsession.version)
	xsession.session.Send(xml)

	xsession.id += 1
}

func (self *XmppServer) sendAuthFeatures(xsession *XmppSession) {
	xml := "<stream:features>" +
		//		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>" +
		"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
		"<mechanism>PLAIN</mechanism>" +
		"</mechanisms>" +
		"</stream:features>"
	xsession.session.Send(xml)
}

func (self *XmppServer) sendBindFeatures(xsession *XmppSession) {
	xml := "<stream:features>" +
		"<bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/>" +
		"</stream:features>"
	xsession.session.Send(xml)
}

func (self *XmppServer) responseAuthPlain(xsession *XmppSession, sa *saslAuth) {
	data, err := base64.StdEncoding.DecodeString(sa.Value)
	if err != nil {
		err1 := fmt.Errorf("Auth information couldn't be decoded as base64. %v", err)
		xsession.session.Panic(err1)
		return
	}

	decodedAuthInfo := strings.Replace(fmt.Sprintf("%q", data), "\"", "", -1)
	authToken := strings.Split(decodedAuthInfo, "\\x00")
	if len(authToken) != 3 {
		err := fmt.Errorf("Invalid auth information.")
		xsession.session.Panic(err)
		return
	}

	id := authToken[1]
	pwd := authToken[2]

	// Authentication
	if id == "test" && pwd == "1234" {
		self.loginUsers.Join(xsession.session)
		xsession.user = id
		xml := "<success xmlns='urn:ietf:params:xml:ns:xmpp-sasl'/>"
		xsession.session.Send(xml)
	} else {
		xml := "<failure xmlns='urn:ietf:params:xml:ns:xmpp-sasl'><not-authorized/></failure>"
		xsession.session.Send(xml)

		err1 := fmt.Errorf("Invalid id or password")
		xsession.session.Panic(err1)
		return
	}

}

func (self *XmppServer) responseBindingResource(xsession *XmppSession, id string) {
	xml := fmt.Sprintf("<iq id='%s' type='result'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'>"+
		"<jid>%s</jid></bind></iq>", id, xsession.Jid())
	xsession.session.Send(xml)
}

func (self *XmppServer) responseIqUnavailable(xsession *XmppSession, iq *clientIQ) {

	xml := fmt.Sprintf("<iq from='%s' to='%s' id='%s' type='error'>"+
		"<error type='cancel'>"+
		"<service-unavailable xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>"+
		"</error>"+
		"</iq>", xsession.host, xsession.Jid(), iq.ID)

	xsession.session.Send(xml)

}

func (self *XmppServer) makeResource() string {
	id, err := uuid.NewV4()
	if err != nil {
		self.Error(err, "uuid.NewV4()")
	}
	return id.String()
}

/* Internal Functions */
///////////////////////////////////////////////////////////////////////////////
