package dhcp

import (
	"errors"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	dhcp "github.com/krolaw/dhcp4"
	"github.com/mistifyio/mistify-agent/client"
)

var NotFound = errors.New("not found")

type Server struct {
	client *client.Client
	config *Config
}

func NewServer(conf *Config) *Server {
	server := &Server{}
	server.config = conf

	c, _ := client.NewClient(&client.Config{Address: conf.Agent})
	server.client = c

	return server
}

func (s *Server) Run() {
	log.WithFields(log.Fields{
		"agent": s.client.Config.Address,
	}).Info("Starting DHCP server")

	conn, err := NewDHCPConnection(s.config.Interfaces)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "dhcp.NewDHCPConnection",
		}).Fatal(err)
	}

	err = dhcp.Serve(conn, s)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "dhcp.Serve",
		}).Fatal(err)
	}
}

func (s *Server) ServeDHCP(packet dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	var replyType dhcp.MessageType
	var logMessage string

	switch msgType {
	case dhcp.Discover:
		replyType = dhcp.Offer
		logMessage = "Discover"
	case dhcp.Request:
		replyType = dhcp.ACK
		logMessage = "Request"
	default:
		return nil
	}

	mac := packet.CHAddr().String()

	log.WithFields(log.Fields{
		"mac": mac,
	}).Info(logMessage)

	nic, err := s.getNic(mac)
	if err != nil {
		log.WithFields(log.Fields{
			"mac":   mac,
			"error": err,
			"func":  "dhcp.Server.getNic",
		}).Info("Couldn't get NIC")
		return nil
	}

	log.WithFields(log.Fields{
		"mac": mac,
		"ip":  nic.Address,
	}).Info("Returning IP")

	replyOpts := dhcp.Options{
		dhcp.OptionRouter:           net.ParseIP(nic.Gateway).To4(),
		dhcp.OptionDomainNameServer: net.IP{8, 8, 8, 8},
		dhcp.OptionSubnetMask:       net.ParseIP(nic.Netmask).To4(),
	}

	reply := dhcp.ReplyPacket(packet, replyType, net.IPv4zero, net.ParseIP(nic.Address).To4(), time.Hour*24*7, replyOpts.SelectOrderOrAll(replyOpts[dhcp.OptionParameterRequestList]))
	return reply
}

func (server *Server) getNic(mac string) (*client.Nic, error) {
	guests, err := server.client.ListGuests()
	if err != nil {
		return nil, err
	}

	for _, guest := range guests {
		for _, nic := range guest.Nics {
			if nic.Mac == mac {
				return &nic, nil
			}
		}
	}
	return nil, NotFound
}
