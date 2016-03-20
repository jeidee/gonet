package chat_server

import (
	"errors"
	"fmt"

	"github.com/jeidee/gonet"
	. "github.com/jeidee/gonet/example/json_protocol/chat_packet"
)

type ChatStub struct {
	gonet.Stub

	server *ChatServer
}

func (stub *ChatStub) ParseAndDo(data *gonet.IncommingData, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			stub.server.Error(err, "[ChatStub:ParseAndDo]")
		}
	}()

	m, ok := data.Data.(map[string]interface{})
	if !ok {
		err = errors.New("Invalid packet")
		return
	}

	id := uint32(m["packet_id"].(float64))

	switch id {
	case PID_REQ_LOGIN:
		stub.OnReqLogin(data.Session, ReqLogin{PacketId: id, Nickname: m["nickname"].(string)})
		break
	case PID_REQ_USER_LIST:
		stub.OnReqUserList(data.Session, ReqUserList{PacketId: id})
		break
	case PID_REQ_SEND_CHAT:
		stub.OnReqSendChat(data.Session, ReqSendChat{PacketId: id, Chat: m["chat"].(string)})
		break
	default:
		err = errors.New("Not supported packet.")
		break
	}
}

func (stub *ChatStub) OnReqLogin(session *gonet.Session, obj ReqLogin) {
	sendPacket := NewResLogin()

	err := stub.server.Login(session, obj.Nickname)

	if err != nil {
		sendPacket.Result = RESULT_FAIL
		stub.server.Debug("Login failed ... %v", err)
	} else {
		sendPacket.Result = RESULT_OK
		stub.server.Info("%v is login.", obj.Nickname)
	}

	stub.server.proxy.ResLogin(session, sendPacket)
}

func (stub *ChatStub) OnReqUserList(session *gonet.Session, obj ReqUserList) {
	stub.server.Debug("OnReqUserList ... %v", obj)
	packet := NewResUserList()
	packet.NumberOfUsers = len(stub.server.chatSessions)
	packet.Users = make([]UserInfo, packet.NumberOfUsers)

	i := 0
	for _, chatSession := range stub.server.chatSessions {
		packet.Users[i] = UserInfo{Nickname: chatSession.nickname}
	}

	stub.server.proxy.ResUserList(session, packet)
}

func (stub *ChatStub) OnReqSendChat(session *gonet.Session, obj ReqSendChat) {
	stub.server.Debug("OnReqSendChat ... %v", obj)

	// broadcast to all
	sender := stub.server.chatSessions[session]
	packet := NewNotifyChat()
	packet.Chat = obj.Chat
	packet.SenderNickname = sender.nickname

	stub.server.lobby.Broadcast(packet)

}
