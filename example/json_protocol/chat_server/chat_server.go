package chat_server

import (
	"errors"

	"github.com/jeidee/gonet"
)

type ChatServer struct {
	gonet.Server

	chatSessions    map[*gonet.Session]*ChatSession
	nicknameIndices map[string]*ChatSession

	lobby *gonet.SessionGroupManager

	lastSessionId int32

	proxy *ChatProxy
}

func NewChatServer(port int16) *ChatServer {
	s := new(ChatServer)

	s.chatSessions = make(map[*gonet.Session]*ChatSession)
	s.nicknameIndices = make(map[string]*ChatSession)
	s.lastSessionId = 0

	stub := new(ChatStub)
	stub.server = s

	s.proxy = new(ChatProxy)
	s.proxy.server = s

	s.Init(port, new(gonet.JsonProtocol), s, stub)

	// This code must place after s.Init(),
	// because SessionGroupManager of server will be create in s.Init().
	lobby, err := s.SessionGroupManager().NewGroup("lobby")
	s.Debug("lobby %v, err %v", lobby, err)
	if err != nil {
		s.Debug("Can't create a new group[%v].", "lobby")
		return nil
	}
	s.lobby = lobby

	return s
}

///////////////////////////////////////////////////////////////////////////////
/* APIs */

func (s *ChatServer) Login(session *gonet.Session, nickname string) error {
	s.Debug("Login ... %v", nickname)

	cs := s.nicknameIndices[nickname]

	s.Info("cs is ...%v", cs)
	if cs == nil {
		s.lastSessionId += 1
		chatSession := NewChatSession(session, s.lastSessionId)
		chatSession.nickname = nickname

		s.chatSessions[session] = chatSession
		s.nicknameIndices[nickname] = chatSession

		s.lobby.Join(session)

		return nil
	} else {
		return errors.New(nickname + " is already logined")
	}
}

func (s *ChatServer) Logout(session *gonet.Session) {
	chatSession := s.chatSessions[session]

	if chatSession != nil {
		s.lastSessionId -= 1
		delete(s.nicknameIndices, chatSession.nickname)
		delete(s.chatSessions, session)

		s.lobby.Leave(session)

		s.Debug("%v is logout.", chatSession.nickname)
	}
}

func (s *ChatServer) GetChatSession(session *gonet.Session) *ChatSession {
	return s.chatSessions[session]
}

/* APIs */
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
/* Server Event Handlers
... This code belong to the goroutine of server's eventLoop. */

func (s *ChatServer) OnAccept(session *gonet.Session) {
	s.Debug("ChatServer ... OnAccept...%v, Cu is %v.", s.lastSessionId, len(s.chatSessions))
}

func (s *ChatServer) OnClose(session *gonet.Session) {
	s.Debug("ChatServer ... OnClose...%v, Cu is %v.", s.lastSessionId, len(s.chatSessions))
	s.Logout(session)
}

func (s *ChatServer) OnIncomming(data *gonet.IncommingData) {
	s.Debug("ChatServer ... OnIncomming...%v", data.Data)

	//	// chat
	//	err := data.Session.Send(data.Data)
	//	if err != nil {
	//		s.Error(err, "Chating failed.")
	//	}

}

/* End of Server Event Handlers */
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
/* Internal Functions */

/* Internal Functions */
///////////////////////////////////////////////////////////////////////////////
