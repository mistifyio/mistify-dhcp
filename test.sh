#!/bin/bash

set -e
set -o pipefail
#set -x

#trap "kill 0" SIGINT SIGTERM EXIT

sudo killall dhclient || true

pushd cmd/test-dhcp
go build
sudo ./test-dhcp &
DHCP_PID=$!
trap "sudo kill $DHCP_PID" SIGINT SIGTERM EXIT
popd


sleep 1

IP=$(sudo dhcping -h de:ad:be:ef:3d:f4 -s 127.0.0.1 -V | awk '/^yiaddr/ {print $NF}' | sort | tail -n 1)

if [ "$IP" != "10.0.0.2" ]; then
    echo "did not get expect ip"
fi



