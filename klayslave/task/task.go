package task

import (
	"math/big"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
)

type Params struct {
	AccGrp   []*account.Account
	Endpoint string
	GasPrice *big.Int
}

type ExtendedTask struct {
	Name    string
	Weight  int
	Fn      func()
	Init    func(params *Params)
	AccGrp  []*account.Account
	EndPint string
}
