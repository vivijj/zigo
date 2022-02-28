// Package babyjub implements the EdDSA over the BabyJubJub curve: ax^2 + y^2 = 1 + d*x^2*y^2
package babyjub

import (
	"fmt"
	"math/big"

	"github.com/vivijj/zigo/crypto/constants"
	"github.com/vivijj/zigo/crypto/ff"
	"github.com/vivijj/zigo/crypto/utils"
)

// babyjub curve:
var (
	// A is one of the babyjub constants.(Coefficient A)
	A *big.Int
	// Aff is A value in *ff.Element representation
	Aff *ff.Element

	// D is one of the babyjub constants.(Coefficient D)
	D *big.Int
	// Dff is D value in *ff.Element representation
	Dff *ff.Element

	// Order of the babyjubjub curve.
	Order *big.Int

	// SubOrder is the order of the subgroup of the babyjub curve that contains the
	// points that we use.
	SubOrder *big.Int

	// B8 is a base point of the babyjub multiplied by 8 to make it a base point of
	// the subgroup in the curve.
	B8 *Point
)

// init initializes global numbers and the subgroup base.
func init() {
	A = utils.NewIntFromString("168700")
	D = utils.NewIntFromString("168696")
	Aff = utils.NewElement().SetBigInt(A)
	Dff = utils.NewElement().SetBigInt(D)

	Order = utils.NewIntFromString(
		"21888242871839275222246405745257275088614511777268538073601725287587578984328",
	)
	SubOrder = new(big.Int).Rsh(Order, 3) // (2^Cofactor)*SubOrder = Order(Cofactor = 3)

	B8 = NewPoint()
	B8.X = utils.NewIntFromString(
		"16540640123574156134436876038791482806971768689494387082833631921987005038935",
	)
	B8.Y = utils.NewIntFromString(
		"20819045374670962167435360035096875258406992893633759881276124905556507972311",
	)
}

// PointProjective is the Point representation in projective coordinates
type PointProjective struct {
	X *ff.Element
	Y *ff.Element
	Z *ff.Element
}

// NewPointProjective creates a new Point in projective coordinates(0,1,1).
// the new point as know as infinity
func NewPointProjective() *PointProjective {
	return &PointProjective{
		X: utils.NewElement().SetZero(),
		Y: utils.NewElement().SetOne(),
		Z: utils.NewElement().SetOne(),
	}
}

// Affine returns the Point from the projective representation
func (p *PointProjective) Affine() *Point {
	if p.Z.Equal(utils.NewElement().SetZero()) {
		return &Point{
			X: big.NewInt(0),
			Y: big.NewInt(0),
		}
	}

	zinv := utils.NewElement().Inverse(p.Z)
	x := utils.NewElement().Mul(p.X, zinv) // X/Z
	y := utils.NewElement().Mul(p.Y, zinv) // Y/Z
	xBig := big.NewInt(0)
	x.ToBigIntRegular(xBig)
	yBig := big.NewInt(0)
	y.ToBigIntRegular(yBig)
	return &Point{
		X: xBig,
		Y: yBig,
	}
}

// Add computes the addition of two points in projective coordinates
// representation
func (p *PointProjective) Add(q *PointProjective, o *PointProjective) *PointProjective {
	// add-2008-bbjlp
	// https://hyperelliptic.org/EFD/g1p/auto-twisted-projective.html#doubling-dbl-2008-bbjlp
	a := utils.NewElement().Mul(q.Z, o.Z)
	b := utils.NewElement().Square(a)
	c := utils.NewElement().Mul(q.X, o.X)
	d := utils.NewElement().Mul(q.Y, o.Y)
	e := utils.NewElement().Mul(Dff, c)
	e.Mul(e, d)
	f := utils.NewElement().Sub(b, e)
	g := utils.NewElement().Add(b, e)
	x1y1 := utils.NewElement().Add(q.X, q.Y)
	x2y2 := utils.NewElement().Add(o.X, o.Y)
	x3 := utils.NewElement().Mul(x1y1, x2y2)
	x3.Sub(x3, c)
	x3.Sub(x3, d)
	x3.Mul(x3, a)
	x3.Mul(x3, f)
	ac := utils.NewElement().Mul(Aff, c)
	y3 := utils.NewElement().Sub(d, ac)
	y3.Mul(y3, a)
	y3.Mul(y3, g)
	z3 := utils.NewElement().Mul(f, g)

	p.X = x3
	p.Y = y3
	p.Z = z3
	return p
}

// Point represents a point of the babyjub curve.
type Point struct {
	X *big.Int
	Y *big.Int
}

// NewPoint creates a new Point（0,1) as infinity point.
func NewPoint() *Point {
	return &Point{X: big.NewInt(0), Y: big.NewInt(1)}
}

// Set copies a Point c into the Point p
func (p *Point) Set(c *Point) *Point {
	p.X.Set(c.X)
	p.Y.Set(c.Y)
	return p
}

// Projective returns a PointProjective from the Point
func (p *Point) Projective() *PointProjective {
	return &PointProjective{
		X: utils.NewElement().SetBigInt(p.X),
		Y: utils.NewElement().SetBigInt(p.Y),
		Z: utils.NewElement().SetOne(),
	}
}

// Mul multiplies the Point q by the scalar s and stores the result in p,
// which is also returned.
func (p *Point) Mul(s *big.Int, q *Point) *Point {
	resProj := NewPointProjective()
	exp := q.Projective()

	for i := 0; i < s.BitLen(); i++ {
		if s.Bit(i) == 1 {
			resProj.Add(resProj, exp)
		}
		exp = exp.Add(exp, exp)
	}
	p = resProj.Affine()
	return p
}

// InCurve returns true when the Point p is in the babyjub curve.
func (p *Point) InCurve() bool {
	x2 := new(big.Int).Set(p.X)
	x2.Mul(x2, x2)
	x2.Mod(x2, constants.Q)

	y2 := new(big.Int).Set(p.Y)
	y2.Mul(y2, y2)
	y2.Mod(y2, constants.Q)

	a := new(big.Int).Mul(A, x2)
	a.Add(a, y2)
	a.Mod(a, constants.Q)

	b := new(big.Int).Set(D)
	b.Mul(b, x2)
	b.Mul(b, y2)
	b.Add(constants.One, b)
	b.Mod(b, constants.Q)

	return a.Cmp(b) == 0
}

// InSubGroup returns true when the Point p is in the subgroup of the babyjub
// curve.
func (p *Point) InSubGroup() bool {
	if !p.InCurve() {
		return false
	}
	res := NewPoint().Mul(SubOrder, p)
	return (res.X.Cmp(constants.Zero) == 0) && (res.Y.Cmp(constants.One) == 0)
}

// PointCoordSign returns the sign of the curve point coordinate.  It returns
// false if the sign is positive and false if the sign is negative.
func PointCoordSign(c *big.Int) bool {
	return c.Cmp(new(big.Int).Rsh(constants.Q, 1)) == 1
}

// PackSignY packs the given sign and the coordinate Y of a point into a 32
// byte array. This method does not check that the values belong to a valid
// Point in the curve.
func PackSignY(sign bool, y *big.Int) [32]byte {
	leBuf := utils.BigIntLEBytes(y)
	if sign {
		leBuf[31] = leBuf[31] | 0x80
	}
	return leBuf
}

// UnpackSignY returns the sign and coordinate Y from a given compressed point.
// This method does not check that the Point belongs to the BabyJubJub curve,
// thus does not return error in such case. This method is intended to obtain
// the sign and the Y coordinate without checking if the point belongs to the
// curve, if the objective is to uncompress a point, Decompress method should
// be used instead.
func UnpackSignY(leBuf [32]byte) (bool, *big.Int) {
	sign := false
	y := big.NewInt(0)
	if (leBuf[31] & 0x80) != 0x00 {
		sign = true
		leBuf[31] = leBuf[31] & 0x7F
	}
	utils.SetBigIntFromLEBytes(y, leBuf[:])
	return sign, y
}

// Compress the point into a 32 byte array that contains the y coordinate in
// little endian and the sign of the x coordinate.
func (p *Point) Compress() [32]byte {
	sign := PointCoordSign(p.X)
	return PackSignY(sign, p.Y)
}

func (p *Point) CompressBi() *big.Int {
	cp := p.Compress()
	bi := new(big.Int)
	return utils.SetBigIntFromLEBytes(bi, cp[:])
}

// Decompress a compressed Point into p, and also returns the decompressed
// Point.  Returns error if the compressed Point is invalid.
func (p *Point) Decompress(leBuf [32]byte) (*Point, error) {
	var sign bool
	sign, p.Y = UnpackSignY(leBuf)
	return PointFromSignAndY(sign, p.Y)
}

// PointFromSignAndY returns a Point from a Sign and the Y coordinate
// x^2 = (1 - y^2) / (a - d*y^2)
func PointFromSignAndY(sign bool, y *big.Int) (*Point, error) {
	var p Point
	p.X = big.NewInt(0)
	p.Y = y
	if p.Y.Cmp(constants.Q) >= 0 {
		return nil, fmt.Errorf("p.y >= Q")
	}

	y2 := new(big.Int).Mul(p.Y, p.Y)
	y2.Mod(y2, constants.Q)
	xa := big.NewInt(1)
	xa.Sub(xa, y2) // xa == 1 - y^2

	xb := new(big.Int).Mul(D, y2)
	xb.Mod(xb, constants.Q)
	xb.Sub(A, xb) // xb = A - d * y^2

	if xb.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("division by 0")
	}
	xb.ModInverse(xb, constants.Q)
	p.X.Mul(xa, xb) // xa / xb
	p.X.Mod(p.X, constants.Q)
	noSqrt := p.X.ModSqrt(p.X, constants.Q)
	if noSqrt == nil {
		return nil, fmt.Errorf("x is not a square mod q")
	}
	if (sign && !PointCoordSign(p.X)) || (!sign && PointCoordSign(p.X)) {
		p.X.Mul(p.X, constants.MinusOne)
	}
	p.X.Mod(p.X, constants.Q)

	return &p, nil
}
