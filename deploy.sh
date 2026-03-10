#!/bin/bash

echo "Updating system..."
apt update -y
apt upgrade -y

echo "Installing Git..."
apt install -y git

echo "Installing Go..."
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

echo "Adding Go to PATH..."
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin

echo "Cloning repo..."
git clone https://github.com/AlexCrowd3/relay-node.git

cd relay-node

echo "Downloading dependencies..."
go mod tidy

echo "Building..."
go build -o relay

echo "Opening firewall..."
apt install -y ufw
ufw allow 8080
ufw --force enable

echo "Starting server..."
nohup ./relay > relay.log 2>&1 &

echo "Relay server started on port 8080"