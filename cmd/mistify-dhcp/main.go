package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-dhcp"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	conf, err := dhcp.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "dhcp.GetConfig",
		}).Fatal(err)
	}

	server := dhcp.NewServer(conf)
	server.Run()
}
