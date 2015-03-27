package dhcp

import (
	"net"

	"code.google.com/p/go.net/ipv4"
	log "github.com/Sirupsen/logrus"
)

// Connection is a dhcp connection. Uses PacketConn and control messages to make
// sure packets are sent over the correct interface.
type Connection struct {
	pconn  *ipv4.PacketConn
	ifaces []int
}

// NewConnection creates a new Connection
func NewConnection(ifaceNames []string) (*Connection, error) {
	var err error
	conn := Connection{}

	addr := net.UDPAddr{IP: net.IPv4zero, Port: 67}
	pconn, err := net.ListenPacket("udp4", addr.String())
	if err != nil {
		return nil, err
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if len(ifaceNames) == 0 {
			// If no interfaces were specified, respond on all of them
			conn.ifaces = append(conn.ifaces, iface.Index)
		} else {
			for _, ifaceName := range ifaceNames {
				if iface.Name == ifaceName {
					conn.ifaces = append(conn.ifaces, iface.Index)
				}
			}
		}
	}

	conn.pconn = ipv4.NewPacketConn(pconn)

	return &conn, nil
}

// ReadFrom reads a packet from the connection
func (conn Connection) ReadFrom(v []byte) (int, net.Addr, error) {
	n, _, src, err := conn.pconn.ReadFrom(v)
	return n, src, err
}

// WriteTo writes a packet to the address
func (conn Connection) WriteTo(data []byte, dest net.Addr) (int, error) {
	var n int
	cm := &ipv4.ControlMessage{}

	destAddr, err := net.ResolveUDPAddr(dest.Network(), dest.String())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "net.ResolveUDPAddr",
		}).Error(err)
		return 0, err
	}

	if net.IPv4bcast.Equal(destAddr.IP) {
		for _, iface := range conn.ifaces {
			cm.IfIndex = iface
			n, err = conn.pconn.WriteTo(data, cm, dest)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"func":  "ipv4.PacketConn.WriteTo",
				}).Error(err)
				return 0, err
			}
		}

		return n, nil
	}
	return conn.pconn.WriteTo(data, cm, dest)
}
