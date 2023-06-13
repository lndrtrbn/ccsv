package main

import (
	"flag"
	"lndr/ccsv/internal/client"
	"log"
)

var (
	ui   *client.UI
	chat *client.Client
)

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
	ui = client.NewUI(*name, chat.Send)

	chat.Start(ui.Append)
	ui.Start()
}
