# WireGuard Control #

Web & API based configuration management for WireGuard VPN.

### Contains ###

* GoLang Server
* WebUI using Vue.js

### Getting Started ###

* Edit .env file for the required setup (See below for Template)
* Install Dependencies - Go & Vue.js
* If wireguard configuration already exists - create the file server.json with required changes
* Allow the required system port on the firewall

### Create Watcher to Restart WireGuard when Configuartion Changes ###

* Create the path file

# /etc/systemd/system/wg-watcher.path
[Unit]
Description=Watch /etc/wireguard for changes

[Path]
PathModified=/etc/wireguard

[Install]
WantedBy=multi-user.target

* Create the service file

# /etc/systemd/system/wg-watcher.service
[Unit]
Description=WireGuard directory watcher
After=network.target

[Service]
Type=oneshot
ExecStart=/bin/systemctl restart wg-quick@wg0.service

[Install]
WantedBy=multi-user.target

Finally Run:
1. sudo systemctl daemon-reload

2. sudo systemctl enable wg-watcher.path && sudo systemctl start wg-watcher.path
>>> Created symlink /etc/systemd/system/multi-user.target.wants/wg-watcher.path → /etc/systemd/system/wg-watcher.path.

3. sudo systemctl enable wg-watcher.service && sudo systemctl start wg-watcher.service
>>> Created symlink /etc/systemd/system/multi-user.target.wants/wg-watcher.service → /etc/systemd/system/wg-watcher.service.

4. sudo systemctl status wg-watcher.path 

5. sudo systemctl status wg-watcher.service

### Template of .env file ###

SERVER=0.0.0.0
PORT=9080
GIN_MODE=debug  #debug or release

WG_CONF_DIR=/etc/wireguard
WG_KEYS_DIR=/etc/wireguard/keys
WG_INTERFACE_NAME=wg0.conf
WG_DOMAIN=<ip-address-or-domain>:<port> #zone.mydomain.com:51820

SMTP_HOST=smtp.email.com
SMTP_PORT=465
SMTP_USERNAME=noreply@email.com
SMTP_PASSWORD="email-password"
SMTP_FROM="My Domain <noreply@email.com>"
