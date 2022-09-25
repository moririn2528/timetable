#!/bin/bash
echo "install"
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu `lsb_release -cs` test"
sudo apt update
sudo apt install -y docker-ce

sudo gpasswd -a $(whoami) docker
sudo chgrp docker /var/run/docker.sock
sudo service docker restart

mkdir timetable/docker/sql-data
mkdir timetable/docker/cloud/log
touch timetable/docker/cloud/log/go.log
touch timetable/docker/cloud/log/go-public.log

cd ~/timetable/api/sql
cat 0_init.sql | sed -e "1 s/timetable/timetable_public/g" | sed -e "2 s/timetable/timetable_public/g" >> 0_init_public.sql
cd ~