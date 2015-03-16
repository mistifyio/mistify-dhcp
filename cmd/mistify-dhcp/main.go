package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-dhcp"
	logx "github.com/mistifyio/mistify-logrus-ext"
)

func main() {
	if err := logx.DefaultSetup("info"); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "logx.DefaultSetup",
		}).Fatal("Could not set up logging")
	}

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
