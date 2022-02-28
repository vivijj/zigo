package account

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestAccountConvert(t *testing.T) {
	a := Account{
		Address: common.HexToAddress("0x2a500A5e1950aea40C22d8885C8DC3c02e99b3E2"),
		PubkeyPair: PubKeyPair{
			X: big.NewInt(1),
			Y: big.NewInt(2),
		},
		Nonce:    1,
		Balances: map[int]*big.Int{1: big.NewInt(1000)},
	}
	start := time.Now()
	for i := 0; i < 1000; i++ {
		a.ToCircuitAcc()
	}
	fmt.Println(time.Since(start))
}
