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
	flag.Parse()

	if *name == "" {
		log.Fatalf("| Error starting the app, please specify a name with --name")
	}

	chat = client.NewClient("http://localhost:3000", *name)
	chat.Start(onMessageFunc)

	ui = client.NewUI(*name, inputDoneFunc)
	ui.Start()
}
