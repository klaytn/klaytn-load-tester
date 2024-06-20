// Package analyticTC implements a blockchain-analytic benchmark similar to BlockBench's Analytic benchmark.
package analyticTC

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sync"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn-load-tester/klayslave/task"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/myzhan/boomer"
)

const Name = "analyticTC"

var (
	endPoint string
	cliPool  clipool.ClientPool

	nAcc   int
	accGrp []*account.Account

	mutex       sync.Mutex
	initialized = false
)

const (
	testQueryTotalTxVal = iota
	testQueryLargestTxVal
	testQueryLargestAccBal
	testLast
)

// Init initializes cliPool and accGrp.
func Init(params *task.Params) {
	mutex.Lock()
	defer mutex.Unlock()

	if !initialized {
		initialized = true

		endPoint = params.Endpoint

		cliCreate := func() interface{} {
			c, err := client.Dial(endPoint)
			if err != nil {
				log.Fatalf("[Analytic] Failed to connect to %s, err=%v", endPoint, err)
			}
			return c
		}

		cliPool.Init(20, 300, cliCreate)

		for _, acc := range params.AccGrp {
			accGrp = append(accGrp, acc)
		}
		nAcc = len(accGrp)
	}
}

func toString(target int) string {
	switch target {
	case testQueryTotalTxVal:
		return "QueryTotalTxVal"
	case testQueryLargestTxVal:
		return "QueryLargestTxVal"
	case testQueryLargestAccBal:
		return "QueryLargestAccBal"
	default:
		return "Unknown"
	}
}

// Run randomly executes one test case.
func Run() {
	target := rand.Int() % testLast
	log.Printf("[Analytic] calling %s()...\n", toString(target))

	switch target {
	case testQueryTotalTxVal:
		QueryTotalTxVal()
	case testQueryLargestTxVal:
		QueryLargestTxVal()
	case testQueryLargestAccBal:
		QueryLargestAccBal()
	default:
	}
}

// QueryTotalTxVal calculates the sum of transaction's values in the latest 30 blocks.
func QueryTotalTxVal() {
	ctx := context.Background()
	conn := cliPool.Alloc().(*client.Client)
	msg := "analytic/QueryTotalTxVal to " + endPoint

	// Get the latest block
	start := boomer.Now()
	block, err := conn.BlockByNumber(ctx, nil)
	if err != nil {
		elapsed := boomer.Now() - start
		log.Printf("[Analytic/QueryTotalTxVal] Failed to call BlockByNumber(), err=%v\n", err)
		boomer.Events.Publish("request_failure", "http", msg, elapsed, err.Error())
		conn.Close()
		return
	}

	blockCnt := 1
	txCnt := len(block.Transactions())
	totalValue := sumUpValues(block.Transactions())

	blockNum := block.Number()

	// Do not continue if the blockchain doesn't have 30 blocks
	if blockNum.Int64() < 30 {
		elapsed := boomer.Now() - start
		log.Printf("[Analytic/QueryTotalTxVal] TC needs 30 blocks, but the blockchain has only %v blocks.\n", blockNum)
		boomer.Events.Publish("request_failure", "http", msg, elapsed, "not enough blocks")
		cliPool.Free(conn)
		return
	}

	// Read 29 more previous blocks from the latest block
	startNum := new(big.Int).Set(blockNum)
	startNum.Sub(startNum, big.NewInt(29))
	if startNum.Cmp(big.NewInt(0)) == -1 {
		startNum = big.NewInt(0)
	}
	for blockNum.Cmp(startNum) > 0 {
		block, err := conn.BlockByNumber(ctx, blockNum)
		if err != nil {
			elapsed := boomer.Now() - start
			log.Printf("[Analytic/QueryTotalTxVal] Failed to call BlockByNumber(%v), err=%v\n", blockNum, err)
			boomer.Events.Publish("request_failure", "http", msg, elapsed, err.Error())
			conn.Close()
			return
		}

		txCnt += len(block.Transactions())
		totalValue.Add(totalValue, sumUpValues(block.Transactions()))
		blockCnt++

		blockNum.Sub(blockNum, big.NewInt(1))
	}
	elapsed := boomer.Now() - start

	fmt.Printf("[Analytic/QueryTotalTxVal] The total value in %d txs from %d latest blocks: %v (%v ms)\n", txCnt, blockCnt, totalValue, elapsed)
	boomer.Events.Publish("request_success", "http", msg, elapsed, int64(10))
	cliPool.Free(conn)
}

func sumUpValues(txs types.Transactions) *big.Int {
	totalValue := big.NewInt(0)
	for _, tx := range txs {
		totalValue.Add(totalValue, tx.Value())
	}
	return totalValue
}

// QueryLargestTxVal finds the largest transaction value in the latest 30 blocks.
func QueryLargestTxVal() {
	ctx := context.Background()
	conn := cliPool.Alloc().(*client.Client)
	msg := "analytic/QueryLargestTxVal to " + endPoint

	// Get the latest block
	start := boomer.Now()
	block, err := conn.BlockByNumber(ctx, nil)
	if err != nil {
		elapsed := boomer.Now() - start
		log.Printf("[Analytic/QueryLargestTxVal] Failed to call BlockByNumber(), err=%v\n", err)
		boomer.Events.Publish("request_failure", "http", msg, elapsed, err.Error())
		conn.Close()
		return
	}

	blockCnt := 1
	txCnt := len(block.Transactions())
	largestValue := findLargestValue(block.Transactions())

	blockNum := block.Number()

	// Do not continue if the blockchain doesn't have 30 blocks
	if blockNum.Int64() < 30 {
		elapsed := boomer.Now() - start
		log.Printf("[Analytic/QueryLargestTxVal] TC needs 30 blocks, but the blockchain has only %v blocks.\n", blockNum)
		boomer.Events.Publish("request_failure", "http", msg, elapsed, "not enough blocks")
		cliPool.Free(conn)
		return
	}

	// Read 29 more previous blocks from the latest block
	startNum := new(big.Int).Set(blockNum)
	startNum.Sub(startNum, big.NewInt(29))
	if startNum.Cmp(big.NewInt(0)) == -1 {
		startNum = big.NewInt(0)
	}
	for blockNum.Cmp(startNum) > 0 {
		block, err := conn.BlockByNumber(ctx, blockNum)
		if err != nil {
			elapsed := boomer.Now() - start
			log.Printf("[Analytic/QueryLargestTxVal] Failed to call BlockByNumber(%v), err=%v\n", blockNum, err)
			boomer.Events.Publish("request_failure", "http", msg, elapsed, err.Error())
			conn.Close()
			return
		}

		txCnt += len(block.Transactions())
		val := findLargestValue(block.Transactions())
		if largestValue.Uint64() < val.Uint64() {
			largestValue.Set(val)
		}
		blockCnt++

		blockNum.Sub(blockNum, big.NewInt(1))
	}
	elapsed := boomer.Now() - start

	fmt.Printf("[Analytic/QueryLargestTxVal] The largest value in %d txs from %d latest blocks: %v (%v ms)\n", txCnt, blockCnt, largestValue, elapsed)
	boomer.Events.Publish("request_success", "http", msg, elapsed, int64(10))
	cliPool.Free(conn)
}

func findLargestValue(txs types.Transactions) *big.Int {
	largestValue := big.NewInt(0)
	for _, tx := range txs {
		if largestValue.Uint64() < tx.Value().Uint64() {
			largestValue.Set(tx.Value())
		}
	}
	return largestValue
}

// QueryLargestAccBal finds the largest balance of a randomly chosen account in the latest 30 blocks.
func QueryLargestAccBal() {
	msg := "analytic/QueryLargestAccBal to " + endPoint
	targetAddr := accGrp[rand.Int()%nAcc].GetAddress()

	ctx := context.Background()
	conn := cliPool.Alloc().(*client.Client)

	// Get the latest block to obtain the block number
	start := boomer.Now()
	block, err := conn.BlockByNumber(ctx, nil)
	if err != nil {
		elapsed := boomer.Now() - start
		log.Printf("[Analytic/QueryLargestAccBal] Failed to call BlockByNumber(), err=%v\n", err)
		boomer.Events.Publish("request_failure", "http", msg, elapsed, err.Error())
		conn.Close()
		return
	}

	blockCnt := 0
	largestBal := big.NewInt(0)

	blockNum := block.Number()

	// Do not continue if the blockchain doesn't have 30 blocks
	if blockNum.Int64() < 30 {
		elapsed := boomer.Now() - start
		log.Printf("[Analytic/QueryLargestAccBal] TC needs 30 blocks, but the blockchain has only %v blocks.\n", blockNum)
		boomer.Events.Publish("request_failure", "http", msg, elapsed, "not enough blocks")
		cliPool.Free(conn)
		return
	}

	// Find targetAddr's largest balance in the 30 latest blocks
	startNum := new(big.Int).Set(blockNum)
	startNum.Sub(startNum, big.NewInt(30))
	if startNum.Cmp(big.NewInt(0)) == -1 {
		startNum = big.NewInt(0)
	}
	for blockNum.Cmp(startNum) > 0 {
		bal, err := conn.BalanceAt(ctx, targetAddr, blockNum)
		if err != nil {
			elapsed := boomer.Now() - start
			log.Printf("[Analytic/QueryLargestAccBal] Failed to call BalanceAt(%v), err=%v\n", blockNum, err)
			boomer.Events.Publish("request_failure", "http", msg, elapsed, err.Error())
			conn.Close()
			return
		}

		if largestBal.Uint64() < bal.Uint64() {
			largestBal.Set(bal)
		}

		blockCnt++
		blockNum.Sub(blockNum, big.NewInt(1))
	}
	elapsed := boomer.Now() - start

	fmt.Printf("[Analytic/QueryLargestAccBal] The largest balance of account %s in %d latest blocks: %v (%v ms)\n", targetAddr.String(), blockCnt, largestBal, elapsed)
	boomer.Events.Publish("request_success", "http", msg, elapsed, int64(10))
	cliPool.Free(conn)
}
