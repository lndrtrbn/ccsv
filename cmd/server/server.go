package main

import (
	"flag"
	"lndr/ccsv/internal/server"
	"log"
)

func main() {
	port := flag.String("port", "", "the port to listen to")
	flag.Parse()

	if *port == "" {
		*port = "3000"
		log.Printf("| No port specified, using default port 3000")
	}

	server.NewServer().ListenAndServe(":" + *port)
}
