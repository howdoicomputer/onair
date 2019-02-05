#!/bin/bash

echo "Starting deployment."

cd /home/pi

rm /usr/local/bin/onair_server
rm /etc/systemd/system/onair.service

cp ./onair_server /usr/local/bin/onair_server
chmod +x /usr/local/bin/onair_server

cp ./scripts/onair.service /etc/systemd/system/onair.service
chmod 664 /etc/systemd/system/onair.service

systemctl daemon-reload
systemctl restart onair.service

rm -rf ./scripts ./onair_server ./onair.tar.gz

systemctl is-active --quiet onair.service

if [ $? -eq 0 ]; then
    echo "Deployment successful."
fi
