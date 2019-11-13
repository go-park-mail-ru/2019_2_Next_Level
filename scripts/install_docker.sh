#!/bin/bash

sudo yum check-update

# Docker

curl -fsSL https://get.docker.com/ | sh

sudo systemctl start docker

sudo systemctl status docker #check the state of docker daemon

sudo systemctl enable docker

sudo usermod -aG docker $(whoami)

# Docker-compose

sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose

sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

docker-compose --version
