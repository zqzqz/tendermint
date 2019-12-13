#!/bin/bash
DIR=$(dirname $0)

nohup /home/ubuntu/go/bin/tendermint node --proxy_app=kvstore --trace 2>&1 > ~/.tendermint/run.log &
echo "started"
timeout 300 bash ${DIR}/send_transaction.sh
pkill tendermint
exit 0
