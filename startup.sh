#!/bin/bash
echo "install"
sudo apt install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu `lsb_release -cs` test"
sudo apt update
sudo apt install -y docker-ce

sudo gpasswd -a $(whoami) docker
sudo chgrp docker /var/run/docker.sock
sudo service docker restart

echo "git clone"
sudo apt install -y git
git clone https://github.com/moririn2528/timetable.git
cd timetable
git checkout moririn
cd ~

mkdir timetable/docker/sql-data
mkdir timetable/docker/aws/log
touch timetable/docker/aws/log/go.log