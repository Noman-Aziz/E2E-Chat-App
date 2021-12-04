package main

import "github.com/Noman-Aziz/E2E-Chat-App/Chat/server"

func main() {
	var s server.ChatServer
	s = server.NewServer()
	s.Listen(":3333")

	// start the server
	s.Start()
}
