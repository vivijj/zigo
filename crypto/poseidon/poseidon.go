// Package poseidon is implementation of the Poseidon permutation(not the latest
// one) Starkad and Poseidon: New Hash Functions for Zero Knowledge Proof
// Systems - Lorenzo Grassi, Daniel Kales, Dmitry Khovratovich, Arnab Roy,
// Christian Rechberger, and Markus Schofnegger
// https://eprint.iacr.org/eprint-bin/getfile.pl?entry=2019/458&version=20190510:122118&file=458.pdf
// This poseidon hash base on BN254(p =
// 21888242871839275222246405745257275088548364400416034343698204186575808495617)
package poseidon

import (
	"math/big"

	"github.com/steinselite/zigo/crypto/internal/ff"
	"github.com/steinselite/zigo/crypto/internal/utils"

	"golang.org/x/crypto/blake2b"
)

type Params struct {
	t         int
	nRoundsF  int
	nRoundsP  int
	constantC []*ff.Element
	constantM [][]*ff.Element
}

// NewParams get the param will use in the poseidon permutation
func NewParams(t int, nRoundsF int, nRoundsP int) Params {
	if t < 2 {
		panic("Invalid param of poseidon")
	}
	constantC := constants(SeedC, nRoundsF+nRoundsP)
	constantM := mdsMatrix(SeedM, t)

	return Params{
		t,
		nRoundsF,
		nRoundsP,
		constantC,
		constantM,
	}
}

// exp5 performs x^5 mod p
func exp5(a *ff.Element) {
	a.Exp(*a, big.NewInt(5))
}

func exp5state(state []*ff.Element) {
	for i := 0; i < len(state); i++ {
		exp5(state[i])
	}
}

func zero() *ff.Element {
	return utils.NewElement()
}

// iacr.org/2019/458 § 2.3 About the MDS Matrix (pg 8)
// Also:
// - https://en.wikipedia.org/wiki/Cauchy_matrix
func mdsMatrix(seed string, t int) [][]*ff.Element {
	cauchyMatrix := constants(seed, t*2)
	m := make([][]*ff.Element, t)
	for i := 0; i < t; i++ {
		m[i] = make([]*ff.Element, t)
		for j := 0; j < t; j++ {
			m[i][j] = zero().Inverse(zero().Sub(cauchyMatrix[i], cauchyMatrix[t+j]))
		}
	}
	return m
}

// computes the constant use in the ARK(Add-Round Key)
func constants(seed string, n int) []*ff.Element {
	res := make([]*ff.Element, n)
	hash := blake2b.Sum256([]byte(seed))

	for i := 0; i < n; i++ {
		newN := utils.SetElementFromLEBytes(zero(), hash[:])
		res[i] = newN
		hash = blake2b.Sum256(hash[:])
	}
	return res
}

// ARK(.) add round key
func ark(state []*ff.Element, conK *ff.Element) {
	for i := 0; i < len(state); i++ {
		state[i].Add(state[i], conK)
	}
}

// § 2.2 The Hades Strategy (pg 6)
// In more details, assume R_F = 2 · R_f is an even number. Then
// - the first R_f rounds have a full S-Box layer,
// - the middle R_P rounds have a partial S-Box layer (i.e., 1 S-Box layer),
// - the last R_f rounds have a full S-Box layer
func sbox(state []*ff.Element, i int, params Params) {
	halfF := params.nRoundsF / 2
	if i < halfF || i >= (halfF+params.nRoundsP) {
		exp5state(state)
	} else {
		exp5(state[0])
	}

}

// the mixing layer is matrix vector product of the state with the mixing matrix
// - https://mathinsight.org/matrix_vector_multiplication
// m is a t*t MDS matrix
func mix(state []*ff.Element, m [][]*ff.Element) []*ff.Element {
	mul := zero()
	newState := make([]*ff.Element, len(m))
	for i := 0; i < len(newState); i++ {
		newState[i] = zero()
	}
	for i := 0; i < len(state); i++ {
		for j := 0; j < len(state); j++ {
			mul.Mul(m[i][j], state[j])
			newState[i].Add(newState[i], mul)
		}
	}
	return newState
}

// Hash calculate the poseidon hash
func Hash(inpbi []*big.Int, params Params) *big.Int {
	state := make([]*ff.Element, params.t)
	input := utils.BigIntArrayToElementArray(inpbi)

	copy(state[:len(input)], input[:])
	for i := len(input); i < len(state); i++ {
		state[i] = zero()
	}

	for i, e := range params.constantC {
		ark(state, e)
		sbox(state, i, params)
		state = mix(state, params.constantM)
	}

	rE := state[0]
	r := big.NewInt(0)
	rE.ToBigIntRegular(r)
	return r
}
