#!/bin/bash

SECURITY_GROUP="sg-daa923a1"
INSTANCE_TYPE="t2.micro"
KEY_NAME="ubuntu2"
IMAGE_ID="ami-0dbf5ea29a7fc7e05"
REGION="us-west-1"

DIR=$(dirname $0)
CMD=$1

if [ "$CMD" = "list" ]; then
    INSTANCE_IDS=${@:2}
    aws ec2 describe-instances --instance-ids $INSTANCE_IDS --query 'Reservations[*].Instances[*].PublicIpAddress' | grep -oE "[0-9]*\.[0-9]*\.[0-9]*\.[0-9]*" > ${DIR}/ips.list
elif [ "$CMD" = "launch" ]; then
    INSTANCE_COUNT=$2
    INSTANCE_IDS=$(aws ec2 run-instances --count $INSTANCE_COUNT --instance-type ${INSTANCE_TYPE} --image-id ${IMAGE_ID}  --key-name ${KEY_NAME} --query 'Instances[*].InstanceId' | grep -oE 'i-[0-9a-z]+')
    echo $INSTANCE_IDS
    for id in $INSTANCE_IDS
    do
      echo "$id" >> ${DIR}/instances.list
    done
    aws ec2 describe-instances --instance-ids $INSTANCE_IDS --query 'Reservations[*].Instances[*].PublicIpAddress' | grep -oE "[0-9]*\.[0-9]*\.[0-9]*\.[0-9]*" > ${DIR}/ips.list
elif [ "$CMD" = "stop" ]; then
    INSTANCE_IDS=${@:2}
    aws ec2 stop-instances --instance-ids ${INSTANCE_IDS}
elif [ "$CMD" = "start" ]; then
    INSTANCE_IDS=${@:2}
    aws ec2 start-instances --instance-ids ${INSTANCE_IDS}
fi
