#!/bin/bash

# Variables
SERVER_USER="ubuntu"
SERVER_IP="ubuntu@ec2-44-221-242-117.compute-1.amazonaws.com"
PROJECT_DIR="/Users/Apple/Documents/code/devops/event-mgtfor this lineF"

# Ensure the project directory exists
ssh -i "event_server.pem" $SERVER_USER@$SERVER_IP "mkdir -p $PROJECT_DIR"

# Build the project
GOOS=linux GOARCH=amd64 go build -o yourapp .

# Copy the binary to the server
scp -i "your-key.pem" yourapp $SERVER_USER@$SERVER_IP:$PROJECT_DIR/yourapp

# Restart the application (adjust based on your setup)
ssh -i "your-key.pem" $SERVER_USER@$SERVER_IP << EOF
  sudo systemctl stop yourapp || true
  sudo mv $PROJECT_DIR/yourapp /usr/local/bin/yourapp
  sudo systemctl start yourapp
EOF
