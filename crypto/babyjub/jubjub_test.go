package babyjub

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vivijj/zigo/crypto/utils"
)

func TestCompressDecompress1(t *testing.T) {
	x := utils.NewIntFromString(
		"17777552123799933955779906779655732241715742912184938656739573121738514868268",
	)
	y := utils.NewIntFromString(
		"2626589144620713026669568689430873010625803728049924121243784502389097019475",
	)
	p := &Point{X: x, Y: y}

	buf := p.Compress()
	fmt.Println(buf)

	assert.Equal(
		t,
		"53b81ed5bffe9545b54016234682e7b2f699bd42a5e9eae27ff4051bc698ce85",
		hex.EncodeToString(buf[:]),
	)

	p2, err := NewPoint().Decompress(buf)
	assert.Equal(t, nil, err)
	assert.Equal(t, p.X.String(), p2.X.String())
	assert.Equal(t, p.Y.String(), p2.Y.String())
}

func TestPackSignY(t *testing.T) {
	x := utils.NewIntFromString(
		"17777552123799933955779906779655732241715742912184938656739573121738514868268",
	)
	y := utils.NewIntFromString(
		"2626589144620713026669568689430873010625803728049924121243784502389097019475",
	)
	p := &Point{
		X: x,
		Y: y,
	}
	res := p.CompressBi()
	fmt.Println("bi is: ", res)
	buf := p.Compress()
	fmt.Println("hex buf is:", hex.EncodeToString(buf[:]))
	bi := new(big.Int)
	bi = utils.SetBigIntFromLEBytes(bi, buf[:])
	fmt.Println("num is :", bi)
}
