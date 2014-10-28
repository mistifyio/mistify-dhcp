Example runit scripts and config.

To install (assuming binaries are in /opt/mistify/sbin):

```
install -d -m 0644 dhcp.json /etc/mistify-dhcp/dhcp.json
install -d -m 0755 run /etc/sv/mistify-dhcp/run
install -d -m 0755 log /etc/sv/mistify-agent/log/run
ln -sf /etc/sv/mistify-dhcp /etc/service
```