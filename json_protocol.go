package gonet

import (
	"encoding/json"
	"net"
)

type JsonProtocol struct {
	encoder *json.Encoder
	decoder *json.Decoder
}

func (self *JsonProtocol) Make(conn net.Conn) Protocol {

	newProtocol := new(JsonProtocol)

	newProtocol.encoder = json.NewEncoder(conn)
	newProtocol.decoder = json.NewDecoder(conn)

	return newProtocol
}

func (self *JsonProtocol) Encode(session *Session, data interface{}) (interface{}, error) {
	err := self.encoder.Encode(data)

	if err != nil {
		return nil, err
	}

	return data, err
}

func (self *JsonProtocol) Decode(session *Session) (interface{}, error) {
	var data interface{}

	err := self.decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
