PREFIX := /opt/mistify
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
	mkdir -p $(DESTDIR)${SBIN_DIR}
	mkdir -p $(DESTDIR)${SV_DIR}

	install -D cmd/mistify-dhcp/mistify-dhcp $(DESTDIR)${SBIN_DIR}/mistify-dhcp
	install -D -m 0755 scripts/sv/run $(DESTDIR)${SV_DIR}/mistify-dhcp/run
	install -D -m 0755 scripts/sv/log $(DESTDIR)${SV_DIR}/mistify-dhcp/log/run


