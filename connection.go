// Connection object. Uses PacketConn and control messages to make sure
// packets are sent over the correct interface.

package dhcp

import (
	"code.google.com/p/go.net/ipv4"
	"github.com/mistifyio/mistify-agent/log"
	"net"
)

type DHCPConnection struct {
	pconn  *ipv4.PacketConn
	ifaces []int
}

func NewDHCPConnection(ifaceNames []string) (*DHCPConnection, error) {
	var err error
	conn := DHCPConnection{}

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

func (conn DHCPConnection) ReadFrom(v []byte) (int, net.Addr, error) {
	n, _, src, err := conn.pconn.ReadFrom(v)
	return n, src, err
}

func (conn DHCPConnection) WriteTo(data []byte, dest net.Addr) (int, error) {
	var n int
	cm := &ipv4.ControlMessage{}

	destAddr, err := net.ResolveUDPAddr(dest.Network(), dest.String())
	if err != nil {
		log.Error(err)
		return 0, err
	}

	if net.IPv4bcast.Equal(destAddr.IP) {
		for _, iface := range conn.ifaces {
			cm.IfIndex = iface
			n, err = conn.pconn.WriteTo(data, cm, dest)
			if err != nil {
				return 0, err
			}
		}

		return n, nil
	} else {
		return conn.pconn.WriteTo(data, cm, dest)
	}
}
