package xmpp_server

import (
	"github.com/jeidee/gonet"
)

const (
	sessionStateNone = iota
	sessionStateTlsNegotiation
	sessionStateSaslNegotiation
	sessionStateResourceBinding
	sessionStateStartSession
)

type XmppSession struct {
	session *gonet.Session

	user string
	host string

	version  string
	id       int32
	resource string

	state int
}

func NewXmppSession(session *gonet.Session) *XmppSession {
	xmppSession := new(XmppSession)
	xmppSession.session = session
	xmppSession.id = 1
	xmppSession.state = sessionStateNone

	return xmppSession
}

///////////////////////////////////////////////////////////////////////////////
/* APIs */

func (self *XmppSession) Jid() string {
	return self.user + "@" + self.host + "/" + self.resource
}

/* APIs */
///////////////////////////////////////////////////////////////////////////////
