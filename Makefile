PREFIX := ${DESTDIR}/opt/mistify
SBIN_DIR=$(PREFIX)/sbin
SV_DIR=$(PREFIX)/sv
ETC_DIR=$(PREFIX)/etc

cmd/mistify-dhcp/mistify-dhcp: cmd/mistify-dhcp/main.go
	cd cmd/mistify-dhcp && \
	go get && \
	go build

clean:
	cd cmd/mistify-agent && \
	go clean

install: cmd/mistify-dhcp/mistify-dhcp
	mkdir -p ${SBIN_DIR}
	mkdir -p ${SV_DIR}

	install -D cmd/mistify-dhcp/mistify-dhcp ${SBIN_DIR}/mistify-dhcp
	install -D -m 0755 scripts/sv/run ${SV_DIR}/mistify-dhcp/run
	install -D -m 0755 scripts/sv/log ${SV_DIR}/mistify-dhcp/log/run

	ln -sf ${SV_DIR}/mistify-dhcp /etc/service/mistify-dhcp


