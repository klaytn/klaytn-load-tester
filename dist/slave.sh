#!/bin/bash

set -u
set -e

MASTERHOST=$1

echo $MASTERHOST
set -v

#./klayslave --max-rps $1 --master-host $2 --master-port 5557 -key $3 -tc=$4  -endpoint $5 > slave.log &
nohup ./klayslave --max-rps $1 --master-host $2 --master-port 5557 -key $3 -tc=$4  -endpoint $5 > slave-$3.log &

#nohup ./klayslave --max-rps 300 --master-host $MASTERHOST --master-port 5557 -key $2 -tc="$3" -endpoint $4 > slave.log &


sleep 2

set +v
