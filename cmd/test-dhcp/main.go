package main

import (
	"encoding/json"
	"github.com/mistifyio/mistify-dhcp"
	"github.com/mistifyio/mistify-dhcp/Godeps/_workspace/src/github.com/mistifyio/mistify-agent/rpc"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func main() {
	guests := rpc.Guests{
		&rpc.Guest{
			Memory: 512,
			Cpu:    1,
			Nics: []rpc.Nic{
				rpc.Nic{
					Model:   "e1000",
					Address: "10.0.0.2",
					Netmask: "255.255.255.0",
					Gateway: "10.0.0.1",
					Mac:     "de:ad:be:ef:3d:f4",
				},
			},
			Disks: []rpc.Disk{
				rpc.Disk{Size: 100},
			},
		},
	}

	guestsJSON, err := json.Marshal(guests)
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(guestsJSON)
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	server := dhcp.NewServer(u.Host, []string{})
	log.Fatal(server.Run())
}
