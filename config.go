package dhcp

import (
	"encoding/json"
	"flag"
	"github.com/mistifyio/mistify-agent/log"
	"io/ioutil"
	"strings"
)

type (
	Config struct {
		Agent      string   `json:"agent"`
		Interfaces []string `json:"interfaces"`
	}
)

func NewConfig() *Config {
	conf := new(Config)
	conf.Agent = "http://127.0.0.1:8080"

	return conf
}

func GetConfig() (*Config, error) {
	conf := new(Config)

	agent := flag.String("agent", "", "Agent address")
	ifaces := flag.String("interfaces", "", "Interfaces to listen on, comma-separated (default: all)")
	configfile := flag.String("config-file", "", "Config file to read (default: none)")

	flag.Parse()

	// Config file is parsed first, other command-line options will override its values
	if *configfile != "" {
		if err := conf.ParseConfigFile(*configfile); err != nil {
			return nil, err
		}
	}

	if *agent != "" {
		conf.Agent = *agent
	}

	if *ifaces != "" {
		conf.Interfaces = strings.Split(*ifaces, ",")
	}

	return conf, nil
}

func (conf *Config) ParseConfigFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		return err
	}

	log.Info("%v\n", conf)

	return nil
}
