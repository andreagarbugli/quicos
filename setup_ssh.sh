#!/bin/bash

# If directory exists, remove it
if [ -d "./configs/ssh" ]; then
    rm -rf ./configs/ssh
fi

# Create directory
mkdir -p ./configs/ssh

CUSTOM_KNOWN_HOSTS_FILE=./configs/ssh/known_hosts
ROUTER1_IP="192.168.50.2"
ROUTER2_IP="192.168.50.3"
CLIENT1_IP="192.168.50.10"
SERVER1_IP="192.168.50.11"

HOSTS=($ROUTER1_IP $ROUTER2_IP $CLIENT1_IP $SERVER1_IP)

# export SSH_KNOWN_HOSTS_FILE=./configs/ssh/known_hosts
export SSH_KNOWN_HOSTS_FILE=$CUSTOM_KNOWN_HOSTS_FILE

for host in "${HOSTS[@]}"; do
    echo "Adding $host to known_hosts"
    ssh-keyscan -H $host >> $CUSTOM_KNOWN_HOSTS_FILE
    echo "Copying SSH key to $host"
    sshpass -p vagrant ssh-copy-id -f -i ~/.ssh/id_rsa.pub -o UserKnownHostsFile=$CUSTOM_KNOWN_HOSTS_FILE vagrant@$host
done