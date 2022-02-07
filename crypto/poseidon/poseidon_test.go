package poseidon

import (
	"fmt"
	"math/big"
	"testing"
)

func TestHash(t *testing.T) {
	param := NewParams(5, 6, 52)
	input := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
	}
	res := Hash(input, param)
	t.Log(res)
	fmt.Println("hash result is: ", res)
}

func BenchmarkHash(b *testing.B) {
	param := NewParams(5, 6, 52)
	var res *big.Int
	for i := 0; i < b.N; i++ {
		input := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
		}
		res = Hash(input, param)
	}
	// b.Log(res)
	fmt.Println("b.n is", b.N)
	fmt.Println("hash result is: ", res)
	fmt.Printf("res is %v and res is %v", res, len(res.Bytes()))
}
