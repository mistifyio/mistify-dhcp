package dhcp

import (
	"bytes"
	"code.google.com/p/go.net/ipv4"
	"crypto/rand"
	"encoding/json"
	"fmt"
	dhcp "github.com/krolaw/dhcp4"
	"github.com/mistifyio/mistify-agent/client"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestDHCP(t *testing.T) {

	guest := client.Guest{
		Memory: 512,
		Cpu:    1,
		Nics: []client.Nic{
			client.Nic{
				Model:   "e1000",
				Address: "10.0.0.2",
				Netmask: "255.255.255.0",
				Gateway: "10.0.0.1",
				Mac:     "DE:AD:BE:EF:3D:F4",
			},
		},
		Disks: []client.Disk{
			client.Disk{Size: 100},
		},
	}

	guestJSON, err := json.Marshal(guest)
	ok(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(guestJSON)
	}))
	defer ts.Close()

	ifname := "eth1"

	// Create UDP connection
	udpconn, err := net.ListenPacket("udp4", net.JoinHostPort(net.IPv4zero.String(), "68"))
	ok(t, err)

	// Create PacketConn with UDP connection as its transport
	pconn := ipv4.NewPacketConn(udpconn)

	// Get interface index and create a control message
	iface, err := net.InterfaceByName(ifname)
	ok(t, err)

	cm := ipv4.ControlMessage{IfIndex: iface.Index}

	// Create a new guest
	resp, err := http.Post(ts.URL+"/guests", "application/json", bytes.NewReader(guestJSON))
	ok(t, err)

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&guest)
	ok(t, err)

	fmt.Printf("Guest: %+v\n", guest)
	// Just need any NIC
	nic := guest.Nics[0]
	fmt.Println(nic.Mac)
	mac, err := net.ParseMAC(nic.Mac)
	ok(t, err)

	// Generate xid for discover packet - 4 random bytes
	xid := make([]byte, 4)
	rand.Reader.Read(xid)

	// Create and send discover packet
	packet := dhcp.RequestPacket(dhcp.Discover, mac, net.IPv4zero, xid, true, []dhcp.Option{})

	_, err = pconn.WriteTo(packet, &cm, &net.UDPAddr{IP: net.IPv4bcast, Port: 67})
	ok(t, err)

	// Listen for a response
	response := make(dhcp.Packet, 512)
	pconn.ReadFrom(response)

	// Make sure it's an offer and the IP address in the response packet matches our guest's
	mtype := dhcp.MessageType(response.ParseOptions()[dhcp.OptionDHCPMessageType][0])
	equals(t, dhcp.Offer, mtype)

	assert(t, response.YIAddr().Equal(net.ParseIP(nic.Address)), "Got wrong IP address: %v != %v", response.YIAddr(), nic.Address)

	// Send back a request packet
	packet = dhcp.RequestPacket(dhcp.Request, mac, response.YIAddr(), xid, true, []dhcp.Option{})

	_, err = pconn.WriteTo(packet, &cm, &net.UDPAddr{IP: net.IPv4bcast, Port: 67})
	ok(t, err)

	// Listen for a response
	pconn.ReadFrom(response)

	// Make sure it's an ack and the IP address in the response packet matches our guest's
	mtype = dhcp.MessageType(response.ParseOptions()[dhcp.OptionDHCPMessageType][0])
	equals(t, dhcp.ACK, mtype)

	assert(t, response.YIAddr().Equal(net.ParseIP(nic.Address)), "Got wrong IP address: %v != %v", response.YIAddr(), nic.Address)
}

// https://github.com/benbjohnson/testing

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
