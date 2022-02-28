// Package constants define some constant big.Int use frequently
package constants

import (
	"fmt"
	"math/big"
)

var (
	// Q is the order of the integer field GF(Q) that fit inside the SNARK
	Q        *big.Int
	Zero     *big.Int
	One      *big.Int
	MinusOne *big.Int
)

// when import this package and init fail,it will panic.
func init() {
	Zero = big.NewInt(0)
	One = big.NewInt(1)
	MinusOne = big.NewInt(-1)

	// This is the order of the finite field
	qString := "21888242871839275222246405745257275088548364400416034343698204186575808495617"
	var ok bool
	Q, ok = new(big.Int).SetString(qString, 10)
	// This should never happen or panic
	if !ok {
		panic(fmt.Sprintf("Bad base 10 string %s", qString))
	}
}
