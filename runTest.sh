#!/bin/bash

MCKEY="3960759054f79e750ecff377444067688903b1234a7f8123c1f77eeff79d69c1"
SCKEY="e3b77cea665fea1bb8547eeabb3aeb2014a1afa12a62aa68e75ab95e1fa99a01"

go build -o runklayslave klayslave/main.go
./runklayslave -max-rps 100 -mcKey=$MCKEY -scKey=$SCKEY -tc="scKLAYTransferTc,scTokenTransferTc"
#./runklayslave -max-rps 1 -mcKey=$MCKEY -scKey=$SCKEY -tc="scTokenTransferTc"
#./runklayslave -max-rps 100 -mcKey=$MCKEY -scKey=$SCKEY -tc="scTokenTransferTc"
