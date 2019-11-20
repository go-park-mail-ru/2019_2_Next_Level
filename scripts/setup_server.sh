#!/bin/bash

# install nslookup, dig, etc
sudo yum install epel-release
sudo yum install -y bind-utils nmap nginx firewalld git

sudo systemctl enable firewalld
sudo systemctl start firewalld
firewall-cmd --zone=public --permanent --add-service=http
firewall-cmd --zone=public --permanent --add-service=https
firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --permanent --add-port=443/tcp
firewall-cmd --permanent --add-port=3001/tcp
firewall-cmd --reload

sudo systemctl start nginx
sudo systemctl enable nginx

sudo useradd go

# Install golang
cd ~ && wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
sudo echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
# Добавление GoPath???
# https://linuxize.com/post/how-to-install-go-on-centos-7/

source /etc/profile

# install postgres
sudo yum install -y postgresql-server postgresql-contrib
sudo postgresql-setup initdb