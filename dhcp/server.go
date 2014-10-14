package dhcp

import (
	dhcp "github.com/krolaw/dhcp4"
	"github.com/mistifyio/mistify-agent/client"
	"github.com/mistifyio/mistify-agent/core"
	"github.com/mistifyio/mistify-agent/log"
	"net"
	"os"
	"time"
)

type Server struct {
	client *client.Client
}

func NewServer(endpoint string) *Server {
	server := &Server{}
	server.client = client.NewClient(endpoint)

	return server
}

func (s *Server) Run() {
	err := dhcp.ListenAndServe(s)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func (s *Server) ServeDHCP(packet dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	var replyType dhcp.MessageType

	switch msgType {
	case dhcp.Discover:
		replyType = dhcp.Offer
	case dhcp.Request:
		replyType = dhcp.ACK
	default:
		return nil
	}

	mac := packet.CHAddr().String()

	log.Info("dhcp.%v: %+v\n", msgType, mac)

	nic, err := s.getNic(mac)
	if err != nil {
		return nil
	}

	log.Info("%+v\n", nic)

	replyOpts := dhcp.Options{
		dhcp.OptionRouter:           net.ParseIP(nic.Gateway).To4(),
		dhcp.OptionDomainNameServer: net.IP{8, 8, 8, 8},
		dhcp.OptionSubnetMask:       net.ParseIP(nic.Netmask).To4(),
	}

	reply := dhcp.ReplyPacket(packet, replyType, net.IP{0, 0, 0, 0}, net.ParseIP(nic.Address).To4(), time.Hour*24*7, replyOpts.SelectOrderOrAll(replyOpts[dhcp.OptionParameterRequestList]))
	return reply
}

func (server *Server) getNic(mac string) (*core.Nic, error) {
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
	return nil, client.NotFound
}

