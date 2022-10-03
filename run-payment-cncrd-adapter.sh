#! /bin/bash

STARTTIME=$(date '+%Y-%m-%d_%H-%M-%S')

echo "Writing log to ~/work/logs/payment-cncrd-adapter.$STARTTIME.log"
echo "Send Ctrl-C/SIGTERM to initiate graceful shutdown"

cd ~/work/payment-cncrd-adapter

./payment-cncrd-adapter -config config.yaml &> ~/work/logs/payment-cncrd-adapter.$STARTTIME.log

