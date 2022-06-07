go build -o runklayslave klayslave/main.go
./runklayslave -key="692bb6bb5ed7008bdfd8ff4d55d460ba628f2c0b9d3afcfa406d2651dd412180" -endpoint="http://127.0.0.1:8545" -tc="transferSignedTx" -vusigned=100000
