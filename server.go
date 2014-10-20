package dhcp

import (
	"errors"
	dhcp "github.com/krolaw/dhcp4"
	"github.com/mistifyio/mistify-agent/client"
	"github.com/mistifyio/mistify-agent/log"
	"net"
	"os"
	"time"
)

var NotFound = errors.New("not found")

type Server struct {
	client     *client.Client
	ifaceNames []string
}

func NewServer(endpoint string, ifaceNames []string) *Server {
	server := &Server{}
	c, _ := client.NewClient(&client.Config{Address: endpoint})
	server.client = c
	server.ifaceNames = ifaceNames

	return server
}

func (s *Server) Run() {
	log.Info("Starting DHCP server, agent address is %s\n", s.client.Config.Address)
	conn, err := NewDHCPConnection(s.ifaceNames)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = dhcp.Serve(conn, s)
	if err != nil {
		log.Error(err)
		os.Exit(1)
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

	log.Info("dhcp.%s: %+v\n", logMessage, mac)

	nic, err := s.getNic(mac)
	if err != nil {
		log.Error("Couldn't get NIC for MAC address %s: %s", mac, err.Error())
		return nil
	}

	log.Info("Returning IP <%s> for MAC address <%s>\n", nic.Address, mac)

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
