package xmpp_server

import (
	"encoding/xml"
)

type Xep interface {
	Parse(xml.StartElement) (xml.Name, interface{}, error)
	Name() string
}
