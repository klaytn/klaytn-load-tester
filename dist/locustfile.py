import json
from web3.auto import w3
from locust import HttpLocust, TaskSet, task

def getver(l):
    url = 'http://localhost:8545'
    payload = {"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":67}
    headers = {'content-type': 'application/json'}
    res = l.client.post('/', data=json.dumps(payload), headers=headers)
    print(res.text)


def createPrivateKey():
    acc = w3.eth.account.create('entropy')
    addr = acc.address
    privateKey = acc.privateKey
    #print('my address is     : {}'.format(acc.address))
    #print('my private key is : {}'.format(acc.privateKey.hex()))
    return acc

def createAccount(l):
    acc = w3.personal.newAccount('1111')
    #print("created new account: " + acc)
    w3.personal.unlockAccount(account=acc, passphrase='1111', duration=0)

    return acc


def sendKlay(l, accFrom, accTo, amount, checkReceipt):
    # myPrivateKey = '0x94cb9f766ef067eb229da85213439cf4cbbcd0dc97ede9479be5ee4b7a93b96f'
    # w3.eth.getTransactionCount(accFrom.address),
    signed_txn = w3.eth.account.signTransaction(dict(
        nonce=0,
        gasPrice=0,
        gas=0x76c0,
        to=accTo.address,
        value=amount,
        data=b'',
    ),
        accFrom.privateKey,
    )

    # print("signed txn=" + format(signed_txn))
    payload = {"jsonrpc": "2.0", "method": "eth_sendRawTransaction", "params": [str(signed_txn.rawTransaction.hex())],
               "id": 1}
    headers = {'content-type': 'application/json'}
    res = l.client.post('/', name="eth_sendTransaction", stream=True, data=json.dumps(payload), headers=headers)
    # print("JSON-RPC Result: " + res.text)
    r = res.json()
    #print("hash:" + r['result'])

    # ret = w3.eth.sendRawTransaction(signed_txn.rawTransaction)
    # print("rawtx=" + str(ret.hex()))

    if checkReceipt:
        receipt = w3.eth.waitForTransactionReceipt(r['result'])
        print(receipt)


def sendSignedTransferTx(l):
    accA = createPrivateKey()
    accB = createPrivateKey()

    for i in range(1000):
        sendKlay(l, accA, accB, 0, 0)


def sendTransferTx(l):
    accA = createAccount(l)
    accB = createAccount(l)
    #hash = w3.eth.sendTransaction({'to': accB, 'from': accA, 'value': 0})

    payload = {"jsonrpc": "2.0", "method": "eth_sendTransaction", "params": [
         {
             "from": accA,
             "to": accB,
             "gas": "0x76c0",
             "gasPrice": "0x0",
             "value": "0x0"
        }
     ], "id": 1}
    headers = {'content-type': 'application/json'}
    res = l.client.post('/', name="eth_sendTransaction", data=json.dumps(payload), headers=headers)
    print("JSON-RPC Result: " + res.text)
    r = res.json()
    print("hash:" + r['result'])


def logout(l):
    print("logout")

    #l.client.post("/logout", {"username":"ellen_key", "password":"education"})

def index(l):
    print("================")
    #l.client.get("/")

def profile(l):
    l.client.post("/profile", {"username":"ellen_key", "password":"education"})
    #l.client.get("/profile")

class UserBehavior(TaskSet):
    #tasks = {sendSignedTransferTx: 1, sendTransferTx: 2}
    #tasks = {sendSignedTransferTx: 1}
    @task(1)
    def pay(self):
        sendSignedTransferTx(self)

    def on_start(self):
        index(self)

    def on_stop(self):
        logout(self)

class WebsiteUser(HttpLocust):
    host = "http://localhost:8546"
    maxsize = 10000

    task_set = UserBehavior
    min_wait = 0
    max_wait = 0
    #min_wait = 5000
    #max_wait = 9000



