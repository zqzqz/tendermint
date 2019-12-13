#!/bin/bash

IP_FILE=$1
CMD=$2

DIR=$(dirname $0)
source ${DIR}/../config.sh

IPS=$(cat $IP_FILE)

echo "running command: { $CMD } on servers..."
for IP in $IPS
do
    ssh -i ${SSH_KEY} ${USERNAME}@${IP} $CMD >> ${DIR}/../log/${IP}.log 2>&1 &
done

wait
