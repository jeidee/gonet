/*
	Ref: https://github.com/mattn/go-xmpself.git
*/
package xmpp_server

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/jeidee/gonet"
)

type tee struct {
	r io.Reader
	w io.Writer
}

func (t tee) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		t.w.Write(p[0:n])
		t.w.Write([]byte("\n"))
	}
	return
}

type XmppProtocol struct {
	encoder *xml.Encoder
	decoder *xml.Decoder
}

func (self *XmppProtocol) Make(conn net.Conn) gonet.Protocol {

	newProtocol := new(XmppProtocol)

	newProtocol.encoder = xml.NewEncoder(conn)
	newProtocol.decoder = xml.NewDecoder(tee{conn, os.Stderr})

	return newProtocol
}

func (self *XmppProtocol) Encode(session *gonet.Session, data interface{}) (interface{}, error) {
	_, err := fmt.Fprintf(session.Conn(), "%v", data)
	return data, err
}

func (self *XmppProtocol) Decode(session *gonet.Session) (interface{}, error) {
	_, data, err := self.next()
	return data, err
}

func (self *XmppProtocol) nextStart() (xml.StartElement, error) {
	for {
		t, err := self.decoder.Token()
		if err != nil && err != io.EOF || t == nil {
			return xml.StartElement{}, err
		}

		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		}
	}
}

func (self *XmppProtocol) next() (xml.StartElement, interface{}, error) {
	// Read start element to find out what type we want.
	se, err := self.nextStart()
	if err != nil {
		return se, nil, err
	}

	// Put it in an interface and allocate one.
	var nv interface{}
	switch se.Name.Space + " " + se.Name.Local {
	case nsStream + " stream":
		return xml.StartElement{}, se, nil
	case nsStream + " features":
		nv = &streamFeatures{}
	case nsStream + " error":
		nv = &streamError{}
	case nsTLS + " starttls":
		nv = &tlsStartTLS{}
	case nsTLS + " proceed":
		nv = &tlsProceed{}
	case nsTLS + " failure":
		nv = &tlsFailure{}
	case nsSASL + " auth":
		nv = &saslAuth{}
	case nsSASL + " mechanisms":
		nv = &saslMechanisms{}
	case nsSASL + " challenge":
		nv = ""
	case nsSASL + " response":
		nv = ""
	case nsSASL + " abort":
		nv = &saslAbort{}
	case nsSASL + " success":
		nv = &saslSuccess{}
	case nsSASL + " failure":
		nv = &saslFailure{}
	case nsBind + " bind":
		nv = &bindBind{}
	case nsClient + " message":
		nv = &clientMessage{}
	case nsClient + " presence":
		nv = &clientPresence{}
	case nsClient + " iq":
		nv = &clientIQ{}
	case nsClient + " error":
		nv = &clientError{}
	default:
		return se, nil, errors.New("unexpected XMPP message " +
			se.Name.Space + " <" + se.Name.Local + "/>")
	}

	if err = self.decoder.DecodeElement(nv, &se); err != nil {
		return se, nil, err
	}

	return se, nv, nil
}
