#! /bin/bash

set -o errexit

if [[ "$RUNTIME_USER" == "" ]]; then
  echo "RUNTIME_USER not set, bailing out. Please run setup.sh first."
  exit 1
fi

mkdir -p tmp
cp payment-cncrd-adapter tmp/
cp config.yaml tmp/
cp run-payment-cncrd-adapter.sh tmp/

chgrp $RUNTIME_USER tmp/*
chmod 640 tmp/config.yaml
chmod 750 tmp/payment-cncrd-adapter
chmod 750 tmp/payment-cncrd-adapter.sh
mv tmp/payment-cncrd-adapter /home/$RUNTIME_USER/work/payment-cncrd-adapter/
mv tmp/config.yaml /home/$RUNTIME_USER/work/payment-cncrd-adapter/
mv tmp/run-payment-cncrd-adapter.sh /home/$RUNTIME_USER/work/
