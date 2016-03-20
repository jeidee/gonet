package chat_client

import (
	//	"errors"

	"github.com/jeidee/gonet"
	. "github.com/jeidee/gonet/example/json_protocol/chat_packet"
)

type ChatClient struct {
	gonet.Client

	nickname string
	isLogin  bool
}

func NewChatClient() *ChatClient {
	c := new(ChatClient)

	stub := new(ChatStub)
	stub.client = c

	c.nickname = ""
	c.isLogin = false

	c.Init(new(gonet.JsonProtocol), c, stub)

	return c
}

///////////////////////////////////////////////////////////////////////////////
/* APIs */

func (c *ChatClient) IsLogin() bool {
	return c.isLogin
}

func (c *ChatClient) ReqLogin(nickname string) {
	obj := NewReqLogin()
	obj.Nickname = nickname

	c.nickname = nickname

	c.Send(obj)
}

func (c *ChatClient) ReqUserList() {
	c.Send(NewReqUserList())
}

func (c *ChatClient) ReqSendChat(chat string) {
	obj := NewReqSendChat()
	obj.Chat = chat

	c.Send(obj)
}

/* APIs */
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
/* Server Event Handlers
... This code belong to the goroutine of server's eventLoop. */

func (c *ChatClient) OnConnect(session *gonet.Session) {
	c.Debug("ChatClient ... OnConnect")
}

func (c *ChatClient) OnClose(session *gonet.Session) {
	c.Debug("ChatClient ... OnClose...")
}

func (c *ChatClient) OnIncomming(data *gonet.IncommingData) {
	c.Debug("ChatClient ... OnIncomming...%v", data.Data)
}

/* End of Server Event Handlers */
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
/* Internal Functions */

/* Internal Functions */
///////////////////////////////////////////////////////////////////////////////
