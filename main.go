package main

import (
	"flag"
	"github.com/mistifyio/mistify-dhcp/dhcp"
)

func main() {
	agent := flag.String("agent", "http://127.0.0.1:8080", "Agent address")
	flag.Parse()

	server := dhcp.NewServer(*agent)
	server.Run()
}
