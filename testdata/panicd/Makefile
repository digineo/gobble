NAME    = gobble-panicd
SERVICE = /etc/systemd/system/$(NAME).service

$(NAME): main.go
	go build -o $@ -ldflags="-s -w" main.go
	go build -o gobble -ldflags="-s -w" ../../*.go

.PHONY: install
install: $(NAME) gobble
	# requires root permissions
	install -m 0755 $+ /usr/local/bin/
	install -m 0644 panicd.service $(SERVICE)
	systemctl daemon-reload
	systemctl restart $(NAME).service

.PHONE: uninstall
uninstall:
	# requires root permissions
	[ -f $(SERVICE) ] && systemctl stop $(NAME).service || :
	rm -f \
		/usr/local/bin/gobble \
		/usr/local/bin/$(NAME) \
		$(SERVICE)
	systemctl daemon-reload

.PHONY: clean
clean:
	rm -f gobble $(NAME)
