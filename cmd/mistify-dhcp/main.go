package main

import (
	"flag"
	"github.com/mistifyio/mistify-dhcp"
	"strings"
)

func main() {
	var ifaceNames []string

	agent := flag.String("agent", "http://127.0.0.1:8080", "Agent address")
	ifaces := flag.String("interfaces", "", "Interfaces to listen on, comma-separated (default: all)")
	flag.Parse()

	if *ifaces != "" {
		ifaceNames = strings.Split(*ifaces, ",")
	}

	server := dhcp.NewServer(*agent, ifaceNames)
	server.Run()
}
