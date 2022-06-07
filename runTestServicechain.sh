#!/bin/bash

MCKEY="76710e519896fb5de792bb439129baf2eb2c4457bca5589adab7ab07a8d2ba3a"
SCKEY="0bbb7187b20bce9ea89c1d4e8ed993d92b9d4c6117c65c98892ac175fce0d9d5"

go build -o runklayslave klayslave/main.go
./runklayslave -max-rps 50 -mcKey=$MCKEY -scKey=$SCKEY -tc="scKLAYTransferTc" > locust.log
#./runklayslave -max-rps 200 -mcKey=$MCKEY -scKey=$SCKEY -tc="scKLAYTransferTc,scTokenTransferTc,scNFTTransferTcWithCheck" > locust.log
#./runklayslave -max-rps 2 -mcKey=$MCKEY -scKey=$SCKEY -tc="scNFTTransferTcWithCheck" > locust.log
#./runklayslave -max-rps 150 -mcKey=$MCKEY -scKey=$SCKEY -tc="scTokenTransferTc" > locust.log
#./runklayslave -max-rps 1 -mcKey=$MCKEY -scKey=$SCKEY -tc="scKLAYTransferTc" > locust.log
