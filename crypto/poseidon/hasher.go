package poseidon

import (
	"fmt"
	"math/big"

	"github.com/vivijj/zigo/crypto/ff"
)

type Hasher struct {
	param Params
}

func NewHasher(t int) *Hasher {
	return &Hasher{
		param: NewParams(t, nRoundsF(t), nRoundsP(t)),
	}
}

func (h *Hasher) HashElements(elements []ff.Element) ff.Element {
	ele := make([]*ff.Element, len(elements))
	// copy the ff.Element here, ude to the HashElement may change the ff.Element if we use *ff.
	// Element
	for i, v := range elements {
		ele[i] = &v
	}
	return *h.HashElementsP(ele)
}

func (h *Hasher) HashElementsP(elements []*ff.Element) *ff.Element {
	return HashElement(elements, h.param)
}
func (h *Hasher) HashBi(inputBi []*big.Int) *big.Int {
	return Hash(inputBi, h.param)
}

// best rounds of f: always 6
func nRoundsF(t int) int {
	return 6
}

func nRoundsP(t int) int {
	switch {
	case t <= 3:
		return 51
	case t > 3 && t <= 7:
		return 52
	case t > 7 && t <= 15:
		return 53
	default:
		// not support when t >= 15 now
		panic(fmt.Sprintf("not support t value: %d", t))
	}
}
