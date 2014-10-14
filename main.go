package main

import (
	"github.com/mistifyio/mistify-dhcp/dhcp"
)

func main() {
	server := dhcp.NewServer("http://127.0.0.1:8080")
	server.Run()
}
