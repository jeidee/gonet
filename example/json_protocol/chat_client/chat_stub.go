package chat_client

import (
	"errors"
	"fmt"

	"github.com/jeidee/gonet"
	. "github.com/jeidee/gonet/example/json_protocol/chat_packet"
)

type ChatStub struct {
	gonet.Stub

	client *ChatClient
}

func (stub *ChatStub) ParseAndDo(data *gonet.IncommingData, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			stub.client.Error(err, "[ChatStub:ParseAndDo]")
		}
	}()

	m, ok := data.Data.(map[string]interface{})
	if !ok {
		err = errors.New("Invalid packet")
		return
	}

	id := uint32(m["packet_id"].(float64))

	switch id {
	case PID_RES_LOGIN:
		stub.OnResLogin(data.Session, ResLogin{PacketId: id, Result: uint16(m["result"].(float64))})
		break
	case PID_RES_USER_LIST:
		stub.OnResUserList(data.Session, ResUserList{PacketId: id})
		break
	case PID_NOTIFY_CHAT:
		stub.OnNotifyChat(data.Session, NotifyChat{PacketId: id, SenderNickname: m["sender_nickname"].(string), Chat: m["chat"].(string)})
		break
	default:
		err = errors.New("Not supported packet.")
		break
	}
}

func (stub *ChatStub) OnResLogin(session *gonet.Session, obj ResLogin) {
	stub.client.Debug("OnResLogin ... %v", obj)

	if obj.Result == RESULT_OK {
		stub.client.isLogin = true
	} else {
		stub.client.isLogin = false
	}

}

func (stub *ChatStub) OnResUserList(session *gonet.Session, obj ResUserList) {
	stub.client.Debug("OnReqUserList ... %v", obj)
}

func (stub *ChatStub) OnNotifyChat(session *gonet.Session, obj NotifyChat) {
	stub.client.Debug("OnNotifyChat ... %v", obj)
}
