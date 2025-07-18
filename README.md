# backend
# GOOS=linux GOARCH=amd64 go build -o MIJKomp
# cd /home/user
# chmod +x MIJKomp
# ./MIJKomp


/etc/systemd/system/MIJKomp.service

[Unit]
Description=MIJKomp Service App
After=network.target

[Service]
User=user
WorkingDirectory=/home/user/MIJKomp/Service
ExecStart=/home/user/MIJKomp/service/MIJKomp
Restart=always

[Install]
WantedBy=multi-user.target


# -----------
sudo systemctl daemon-reexec
sudo systemctl enable MIJKomp
sudo systemctl start MIJKomp
sudo systemctl status MIJKomp

dev
321Admin!@#
51.79.255.146