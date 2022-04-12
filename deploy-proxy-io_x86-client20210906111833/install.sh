#!/bin/bash -il

CLIENT_PATH="$(cd "$(dirname "$0")" && pwd)"
echo $CLIENT_PATH

SVC_PATH=/etc/systemd/system/mimiclet.service

cat > $SVC_PATH << EOF
[Unit]
Description=mimiclet
After=network-online.target systemd-networkd.service systemd-resolved.service
Requires=network-online.target

[Service]
Type=forking
Restart=always
RestartSec=1
ExecStart=$CLIENT_PATH/start.sh
ExecStop=$CLIENT_PATH/stop.sh

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable mimiclet
systemctl start mimiclet
systemctl status mimiclet

if [ $? -eq 0 ]
then
    echo "mimiclet installed!"
else
    echo "mimiclet install failed!"
fi
