package main

import (
	"github.com/jeidee/gonet/example/xmpp/xmpp_server"
)

func main() {

	server := xmpp_server.NewXmppServer(5223)

	if !server.Run() {
		server.Info("Starting server failed!")
	}

	server.Info("Bye!")
}
