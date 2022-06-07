#!/bin/bash

set -u
set -e

set -v
nohup locust -f locustfile.py --master > master.log &
#locust -f locustfile.py --master --host=http://localhost:8545
sleep 2

set +v

