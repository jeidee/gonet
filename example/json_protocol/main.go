package main

import (
	"time"

	"github.com/jeidee/gonet/example/json_protocol/chat_client"
	"github.com/jeidee/gonet/example/json_protocol/chat_server"
)

func main() {

	server := chat_server.NewChatServer(12345)

	// test protocol by client
	go func() {

		for {
			if server.IsRunning() {
				client := chat_client.NewChatClient()

				err := client.Connect("localhost", 12345)
				if err != nil {
					client.Debug("Connect to server failed.")
					time.Sleep(3 * time.Second)
					continue
				}

				// 1. login
				client.ReqLogin("hong")

				// 2. get user list
				for {
					client.Debug("IsLogin ... %v", client.IsLogin())
					if client.IsLogin() {
						client.ReqUserList()
						break
					}
					time.Sleep(5 * time.Second)
				}

				// 3. send caht
				client.ReqSendChat("hello, world!")

				time.Sleep(1 * time.Second)
				client.Close()
			}
			time.Sleep(3 * time.Second)
		}
	}()

	if !server.Run() {
		server.Info("Starting server failed!")
	}

	server.Info("Bye!")
}
