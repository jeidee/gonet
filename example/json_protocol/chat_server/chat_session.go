package chat_server

import (
	"github.com/jeidee/gonet"
)

type ChatSession struct {
	session *gonet.Session

	id       int32
	nickname string
}

func NewChatSession(session *gonet.Session, id int32) *ChatSession {
	chatSession := new(ChatSession)
	chatSession.session = session
	chatSession.id = id

	return chatSession
}
