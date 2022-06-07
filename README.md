# locust-load-tester for service chain load test
(https://github.com/ground-x/locust-load-tester)

## Links

* Locust Website: <a href="http://locust.io">locust.io</a>
* Locust Documentation: <a href="http://docs.locust.io">docs.locust.io</a>
* [Confluence Link](https://groundx.atlassian.net/wiki/spaces/TEC/pages/182091965/locust)

## Description

locust-load-tester is forked from Boomer.
It provides built-in test cases that run to the klaytn node.

It is a better load generator for locust, written in golang. It can spawn thousands of goroutines to run your code concurrently.
It will listen and report to the locust master automatically, your test results will be displayed on the master's web UI.

## Prerequisite
Now, locust 0.10.0 is not compatable with this locust-load-tester. Please install locust 0.9.0.
### how to install locust master
```bash
$ brew install libev
$ pip install locustio==0.9.0
```

## Build

```bash
$ git clone git@github.com:ground-x/locust-load-tester.git
$ cd locust-load-tester/klayslave
$ go get github.com/asaskevich/EventBus
$ go get github.com/ugorji/go/codec
$ go get github.com/zeromq/gomq
$ go get github.com/zeromq/gomq/zmtp
$ go build -o runklayslave klayslave/main.go
```

## Run

### Run Locust master
Locust master can be run like below command. You can use `dummy.py`.

```bash
$ locust -f dummy.py --master
```

### Run locust-load-tester

```bash
$ ./klayslave -max-rps 10 -mcKey="74d94a471b8c88eab6f8cd4322047cd4bccd54dc93b861d2509d158a026eeec7" -scKey="e70e12c0be35c27a157a66d4f1ae84ab1a25d4673f63af6f081284ccae6f942b" -tc="scValueTransferTc"
```
#### Parameters
* --max-rps: limit of request per second
* --mcKey: main chain private key to fund to internal klay test accounts that created before run test case.
* --scKey: service chain private key to fund to internal klay test accounts that created before run test case.
* --numMcUser : number of accounts for signed transaction to use in test case. (default: 10)
* --numScUser: number of accounts for unsigned transaction to use in test case. (default: 10)
* --mcEndpoint: klay node rpc main chain endpoint (default: http://localhost:8545).
* --scEndpoint: klay node rpc service chain endpoint (defalut: http://localhost:7545).

* --master-host: master host ip address (optional)
* --master-port: master port number (optional)



### Troubleshooting


#### Slave가 연결이 안되는 경우
  * klaytn-node RPC 확인
    * klaytn-node에 대한 RPC 연결 확인
      * telnet으로 확인 가능
    * RPC가 Enable 되어 있는지 확인
      * --rpcapi "db,txpool,klay,net,web3,miner,personal,admin,rpc"
  * klaytn-node가 블록을 생산하고 있는지 확인
  * Slave AWS 인스턴스에 접속하여 klayslave의 로그 확인
      * ~/locust 폴더에 slave-(private key)  형식으로 로그파일을 생성하고 있음

    <pre><code>
    $ ssh -i ~/asw/klay-load.pem ubuntu@13.209.26.32


    ubuntu@ip-172-31-20-184:~$ cd locust
    ubuntu@ip-172-31-20-184:~/locust$ ls
    klayslave  locustfile.py  master.log  master.sh  __pycache__  slave-349343aad78f398528907e62b62ce7e7e3c9f57c674e12bbd03857682353a73f.log  slave.log  slave.sh
    ubuntu@ip-172-31-20-184:~/locust$


    ubuntu@ip-172-31-20-184:~/locust$ vi slave-349343aad78f398528907e62b62ce7e7e3c9f57c674e12bbd03857682353a73f.log
    </code></pre>
#### 상기 klayslave 로그에서 Transaction 등 오류가 나는 경우(오류 메시지에 따라서 대응)

  * HTTP 연결이 실패한 경우 : slave와 node가 서로 연결가능한 네트워크인지 확인
    * e.g. Failed to connect RPC ...
  * 제공된 private key 값이 올바르지 않은 경우
    * e.g. Failed to HexToECDSA, Failed to Unlock
  * 제공된 Private key에 대응되는 address 대해서 balance 가 충분한 address인지 확인(대상 노드의 console에서 직접 확인)
  * RPC에 없는 Method라고 출력되는 경우
    * 대상 노드에 RPC가 enable되어 있는지 확인

### Locust의 User, Hatch rate, RPS의 정확한 정의 ###

  * **Number of User**: 모든 slave가 생성하는 GoRoutine 개수의 합, 설정한 값을 각 slave가 균등하게 할당받는다. 즉, user가 10이고 slave가 2개인 경우, 전체 GoRoutine은 10개가 생성되고, 각 slave가 5개씩 생성한다.
  * **Hatch Rate** (users hatched/second) : 목표 GoRoutine 개수까지 도달하는 속도, 즉 위에서 설정한 number of user 값까지 도달하는 속도. number of user=100, hatch rate=10인 경우, 10(=100/10)초 동안 GoRoutine이 생성된다.
  * **RPS** (Request per second) : 하나의 slave가 자신의 endpoint에 요청하는 부하의 정도를 나타낸다. 각 slave를 실행하면 위에서 할당받은 GoRoutine으로 해당 RPS를 만들어 낸다. (slave 실행시 입력해주는 --max-rps의 값은 최대로 만들어 내력 노력하는 목표치이다.) RPS=400이고, slave=4인 경우 전체 네트워크가 받는 요청의 개수는 1,600.
  * 참고 코드
      runner.go
  ```
  func (r *runner) spawnGoRoutines(spawnCount int, quit chan bool) {

    log.Println("Hatching and swarming", spawnCount, "clients at the rate", r.hatchRate, "clients/s...")

    weightSum := 0
    for _, task := range r.tasks {
      weightSum += task.Weight
    }

    for _, task := range r.tasks {

      percent := float64(task.Weight) / float64(weightSum)
      amount := int(round(float64(spawnCount)*percent, .5, 0))

      if weightSum == 0 {
        amount = int(float64(spawnCount) / float64(len(r.tasks)))
      }

      for i := 1; i <= amount; i++ {
        select {
        case <-quit:
          // quit hatching goroutine
          return
        default:
          if i%r.hatchRate == 0 {
            time.Sleep(1 * time.Second)
          }
          atomic.AddInt32(&r.numClients, 1)
          go func(fn func()) {
            for {
              select {
              case <-quit:
                return
              default:
                if rpsControlEnabled {
                  token := atomic.AddInt64(&rpsThreshold, -1)
                  if token < 0 {
                    // RPS threshold is reached, wait until next second
                    <-rpsControlChannel
                  } else {
                    r.safeRun(fn)
                  }
                } else {
                  r.safeRun(fn)
                }
              }
            }
          }(task.Fn)
        }

      }

    }

    r.hatchComplete()

  }
  ```

## License

Open source licensed under the MIT license (see _LICENSE_ file for details).

