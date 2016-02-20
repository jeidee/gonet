package main

import (
	"github.com/jeidee/gonet"
)

type EchoServer struct {
}

func main() {

	server := gonet.NewServer(12345, new(gonet.StringProtocol))
	if !server.Run() {
		server.Info("Starting server failed!")
	}

	server.Info("Bye!")
}
