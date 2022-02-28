package transaction

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/vivijj/zigo/crypto/babyjub"
)

func TestPubkeyUpdateTx_IsAuthDataValid(t *testing.T) {
	x, _ := new(big.Int).SetString(
		"6015732762062863111074994669157831572442835354233667138425543886338497301436", 10,
	)
	y, _ := new(big.Int).SetString(
		"18709592006753622366770445439863348377039354487575630761498799759906802783496", 10,
	)
	authdata, err := hex.DecodeString(
		"9feb33446ae57b7c73472ea7031c2c35080de949d134bbcb241ad16194fee6ce33b86b585b106391bc2f8c12b040483e3d126be16a54240d9b292f160c82be271b",
	)
	if err != nil {
		fmt.Println(err)
	}
	tx := PubkeyUpdateTx{
		AccountId:  1,
		Nonce:      0,
		ValidUntil: 3392838427,
		FeeToken:   0,
		Fee:        big.NewInt(0),
		Account:    common.HexToAddress("0xa749cdefd2d9590549df709bbffec04a9bd35b42"),
		PubKey: babyjub.PublicKey{
			X: x,
			Y: y,
		},
		AuthData: authdata,
	}
	ok := tx.IsAuthDataValid()
	assert.Equal(t, true, ok)
}
