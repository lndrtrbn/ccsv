package main

import (
	"flag"
	"lndr/ccsv/internal/client"
	"log"
	"strings"
)

var (
	ui   *client.UI
	chat *client.Client
)

func inputDoneFunc(value string) {
	chat.Send(value)
}

func onMessageFunc(msg client.Message) {
	pad := 15 - len(msg.Name)
	if pad < 0 {
		pad = 0
	}

	prefix := msg.Name + strings.Repeat(" ", pad) + "| "
	ui.Append(prefix + msg.Content)
}

func main() {
	name := flag.String("name", "", "your pseudo in chat")
	server := flag.String("server", "", "address of the HTTP server")
	flag.Parse()

	if *name == "" {
		log.Fatalf("| Error starting the app, please specify a name with --name")
	}

	if *server == "" {
		*server = "http://localhost:3000"
	}

	chat = client.NewClient(*server, *name)
	chat.Start(onMessageFunc)

	ui = client.NewUI(*name, inputDoneFunc)
	ui.Start()
}
