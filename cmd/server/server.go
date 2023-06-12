package main

import "lndr/ccsv/internal/server"

func main() {
	server.NewServer().ListenAndServe(":3000")
}
