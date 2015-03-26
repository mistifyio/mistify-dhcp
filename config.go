package dhcp

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type (
	// Config contains server configuration
	Config struct {
		Agent      string   `json:"agent"`
		Interfaces []string `json:"interfaces"`
	}
)

// NewConfig creates a new Config
func NewConfig() *Config {
	conf := new(Config)
	conf.Agent = "http://127.0.0.1:8080"

	return conf
}

// GetConfig loads configuration from a config file and flags
func GetConfig() (*Config, error) {
	conf := new(Config)

	var agent, ifaces, configfile string

	flag.StringVarP(&agent, "agent", "a", "", "Agent address")
	flag.StringVarP(&ifaces, "interfaces", "i", "", "Interfaces to listen on, comma-separated (default: all)")
	flag.StringVarP(&configfile, "config-file", "c", "", "Config file to read (default: none)")
	flag.Parse()

	// Config file is parsed first, other command-line options will override its values
	if configfile != "" {
		if err := conf.ParseConfigFile(configfile); err != nil {
			return nil, err
		}
	}

	if agent != "" {
		conf.Agent = agent
	}

	if ifaces != "" {
		conf.Interfaces = strings.Split(ifaces, ",")
	}

	return conf, nil
}

// ParseConfigFile reads a config file
func (conf *Config) ParseConfigFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"config": conf,
	}).Info("Config file parsed")

	return nil
}
