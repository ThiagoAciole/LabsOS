#!/bin/sh
set -eu

[ "$(id -u)" = 0 ] || { echo "execute como root" >&2; exit 1; }

install -m 0755 /home/agent/labs-api-linux /usr/local/bin/labs-api
install -m 0755 /home/agent/labsd-linux /usr/local/bin/labsd
install -m 0755 /home/agent/labs-dashboard-linux /usr/local/bin/labs-dashboard
install -m 0644 /home/agent/labs-api.service /etc/systemd/system/labs-api.service
install -m 0644 /home/agent/labsd.service /etc/systemd/system/labsd.service
install -m 0644 /home/agent/labs-dashboard.service /etc/systemd/system/labs-dashboard.service
install -d -o root -g root -m 0755 /opt/labsos/dashboard
cp -a /home/agent/dashboard/. /opt/labsos/dashboard/

install -d -o root -g data-admin -m 2770 /DATA/Apps/jellyfin
install -d -o root -g data-admin -m 2770 /DATA/Apps/syncthing /DATA/Sync
install -d -o agent -g agent -m 0750 /var/lib/labsos/catalog
install -d -o root -g root -m 0755 /opt/labsos/apps/jellyfin
install -m 0644 /home/agent/jellyfin.compose.yml /opt/labsos/apps/jellyfin/compose.yaml
install -d -o root -g root -m 0755 /opt/labsos/apps/syncthing
install -m 0644 /home/agent/syncthing.compose.yml /opt/labsos/apps/syncthing/compose.yaml

docker compose -f /opt/labsos/apps/jellyfin/compose.yaml config --quiet
docker compose -f /opt/labsos/apps/syncthing/compose.yaml config --quiet
systemctl daemon-reload
systemctl enable labsd.service
systemctl restart labsd.service
systemctl daemon-reload
systemctl enable labs-dashboard.service
systemctl restart labs-dashboard.service
systemctl restart labs-api.service

systemctl is-active labsd.service
systemctl is-active labs-api.service
stat -c '%A %U:%G %n' /run/labsos/labsd.sock
