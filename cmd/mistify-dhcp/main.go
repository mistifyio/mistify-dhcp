package main

import (
	"github.com/mistifyio/mistify-agent/log"
	"github.com/mistifyio/mistify-dhcp"
	"os"
)

func main() {
	conf, err := dhcp.GetConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	server := dhcp.NewServer(conf)
	server.Run()
}
