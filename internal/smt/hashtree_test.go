package smt

import (
	"fmt"
	"testing"
	"time"

	"github.com/vivijj/zigo/crypto/ff"
)

func TestHashTree(t *testing.T) {
	ht := NewHashTree(16, ff.NewElement(12345677))
	start := time.Now()
	for i := 0; i < 1000; i++ {
		ht.Update(1, ff.NewElement(uint64(i)))
	}
	fmt.Println(time.Since(start))
	// path := ht.MerklePath(1)
	// fmt.Println(path)
	// res := ht.VerifyProof(path, 1, ff.NewElement(9999999))
	// fmt.Println(res)
}

func BenchmarkHashTree(b *testing.B) {
	ht := NewHashTree(16, ff.NewElement(12345677))
	ht.Update(1, ff.NewElement(9999999))
	path := ht.MerklePath(1)
	ht.VerifyProof(path, 1, ff.NewElement(9999999))
}
