package gonet

type ServerEventHandler interface {
	OnAccept(session *Session)
	OnClose(session *Session)
	OnIncomming(data *IncommingData)
}

type ClientEventHandler interface {
	OnConnect(session *Session)
	OnClose(session *Session)
	OnIncomming(data *IncommingData)
}
