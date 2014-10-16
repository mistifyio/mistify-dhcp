package main_test

import (
	"bytes"
	"code.google.com/p/go.net/ipv4"
	"crypto/rand"
	"encoding/json"
	dhcp "github.com/krolaw/dhcp4"
	"github.com/mistifyio/mistify-agent/rpc"
	"net"
	"net/http"
	"testing"
)

func TestDHCP(t *testing.T) {
	agent := "http://10.233.22.100:8080"
	ifname := "eth1"

	guest := rpc.Guest{
		Memory: 512,
		Cpu:    1,
		Nics: []rpc.Nic{
			rpc.Nic{
				Model:   "e1000",
				Address: "10.0.0.2",
				Netmask: "255.255.255.0",
				Gateway: "10.0.0.1",
			},
		},
		Disks: []rpc.Disk{
			rpc.Disk{Size: 100},
		},
	}

	guestJSON, err := json.Marshal(guest)
	if err != nil {
		t.Fatalf("Couldn't decode JSON: %s\n", err.Error())
	}

	// Create UDP connection
	udpconn, err := net.ListenPacket("udp4", net.JoinHostPort(net.IPv4zero.String(), "68"))
	if err != nil {
		t.Fatalf("ListenPacket failed: %s\n", err.Error())
	}

	// Create PacketConn with UDP connection as its transport
	pconn := ipv4.NewPacketConn(udpconn)

	// Get interface index and create a control message
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		t.Fatalf("Can't get interface: %s\n", err.Error())
	}

	cm := ipv4.ControlMessage{IfIndex: iface.Index}

	// Create a new guest
	resp, err := http.Post(agent+"/guests", "application/json", bytes.NewReader(guestJSON))
	if err != nil {
		t.Fatalf("Couldn't POST: %s\n", err.Error())
	}

	defer resp.Body.Close()

	// Read the response into the "guest" object created earlier
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		t.Fatalf("Can't read response: %s\n", err.Error())
	}

	json.Unmarshal(buf.Bytes(), &guest)

	// Just need any NIC
	nic := guest.Nics[0]
	mac, err := net.ParseMAC(nic.Mac)
	if err != nil {
		t.Fatalf("Can't parse MAC: %s\n", err.Error())
	}

	// Generate xid for discover packet - 4 random bytes
	xid := make([]byte, 4)
	rand.Reader.Read(xid)

	// Create and send discover packet
	packet := dhcp.RequestPacket(dhcp.Discover, mac, net.IPv4zero, xid, true, []dhcp.Option{})

	_, err = pconn.WriteTo(packet, &cm, &net.UDPAddr{IP: net.IPv4bcast, Port: 67})
	if err != nil {
		t.Fatalf("Can't write packet: %s\n", err.Error())
	}

	// Listen for a response
	response := make(dhcp.Packet, 512)
	pconn.ReadFrom(response)

	// Make sure it's an offer and the IP address in the response packet matches our guest's
	mtype := dhcp.MessageType(response.ParseOptions()[dhcp.OptionDHCPMessageType][0])
	if mtype != dhcp.Offer {
		t.Fatalf("Response was not an offer\n")
	}

	if !response.YIAddr().Equal(net.ParseIP(nic.Address)) {
		t.Fatalf("Got wrong IP address (%v != %v)\n", response.YIAddr(), nic.Address)
	}

	// Send back a request packet
	packet = dhcp.RequestPacket(dhcp.Request, mac, response.YIAddr(), xid, true, []dhcp.Option{})

	_, err = pconn.WriteTo(packet, &cm, &net.UDPAddr{IP: net.IPv4bcast, Port: 67})
	if err != nil {
		t.Fatalf("Can't write packet: %s\n", err.Error())
	}

	// Listen for a response
	pconn.ReadFrom(response)

	// Make sure it's an ack and the IP address in the response packet matches our guest's
	mtype = dhcp.MessageType(response.ParseOptions()[dhcp.OptionDHCPMessageType][0])
	if mtype != dhcp.ACK {
		t.Fatalf("Response was not an ACK\n")
	}

	if !response.YIAddr().Equal(net.ParseIP(nic.Address)) {
		t.Fatalf("Got wrong IP address (%v != %v)\n", response.YIAddr(), nic.Address)
	}
}
