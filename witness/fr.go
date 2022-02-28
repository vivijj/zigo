package witness

import (
	"encoding/json"
	"math/big"

	"github.com/vivijj/zigo/crypto/ff"
)

// Fr is the format that will compatible with the c++ code
// it's a simple wrapper of ff.Element to provide more about mas
type Fr ff.Element

func (f Fr) MarshalJSON() ([]byte, error) {
	ffe := ff.Element(f)
	fstr := ffe.String()
	return json.Marshal(fstr)
}
func (f Fr) UnmarshalJSON(bytes []byte) error {
	var fstr string
	err := json.Unmarshal(bytes, &fstr)
	if err != nil {
		return err
	}
}

func BigIntToFr(bi *big.Int) *Fr {
	return (*Fr)(new(ff.Element).SetBigInt(bi))
}
