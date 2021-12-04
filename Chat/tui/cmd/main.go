package main

import (
	"flag"
	"log"
	"os"

	"github.com/Noman-Aziz/E2E-Chat-App/Chat/client"
	"github.com/Noman-Aziz/E2E-Chat-App/Chat/tui"
)

func main() {
	address := flag.String("server", "", "Which server to connect to")
	first := flag.Bool("first", false, "Are you the first client")
	help := flag.Bool("help", false, "Display Help Page")

	flag.Parse()

	if *help || len(os.Args) < 2 {
		flag.Usage()
		return
	}

	client := client.NewClient(*first)
	err := client.Dial(*address)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	// start the client to listen for incoming message
	go client.Start()

	tui.StartUi(client)
}
