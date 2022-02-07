package babyjub

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignVerifyPoseidon(t *testing.T) {
	// test if the Sign poseidon result is same
	var k PrivateKey
	_, err := hex.Decode(k[:],
		[]byte("0001020304050607080900010203040506070809000102030405060708090001"))
	require.Nil(t, err)
	msgBuf, _ := hex.DecodeString("00010203040506070809")
	msg := new(big.Int).SetBytes(msgBuf[:])

	pk := k.Public()

	assert.Equal(t,
		"1977352365943245022996253601232159691338502210969046081776117504959451086174",
		pk.X.String())
	assert.Equal(
		t,
		"10543573034607953434502791003636731705904014543307411222488284830065095079344",
		pk.Y.String(),
	)
	sig := k.SignPoseidon(msg)
	assert.Equal(t,
		"20714789646677110885733157001561723109447555472354270732814523613007921161743",
		sig.R8.X.String())
	assert.Equal(t,
		"465385463534013514723282026050202550230109697427749352279562400188174024494",
		sig.R8.Y.String())
	assert.Equal(t,
		"964185524797934235033519387723298708154528924757443749906932433779618932229",
		sig.S.String())
	ok := pk.VerifyPoseidon(msg, sig)
	assert.Equal(t, true, ok)
}
