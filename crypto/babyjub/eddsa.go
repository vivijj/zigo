// Package babyjub eddsa implements the EdDSA over the BabyJubJub curve
package babyjub

import (
	"math/big"

	"github.com/steinselite/zigo/crypto/internal/utils"
	"github.com/steinselite/zigo/crypto/poseidon"
)

// PrivateKey is an EdDSA private key, which is a 32byte buffer.
type PrivateKey [32]byte

// Scalar converts a private key into the scalar value s
func (k *PrivateKey) Scalar() *PrivKeyScalar {
	s := new(big.Int).SetBytes(k[:])
	return NewPrivKeyScalar(s)
}

// Public returns the public key corresponding to a private key.
func (k *PrivateKey) Public() *PublicKey {
	return k.Scalar().Public()
}

// PrivKeyScalar represents the scalar s output of a private key
type PrivKeyScalar big.Int

// NewPrivKeyScalar creates a new PrivKeyScalar from a big.Int
func NewPrivKeyScalar(s *big.Int) *PrivKeyScalar {
	sk := PrivKeyScalar(*s)
	return &sk
}

// Public returns the public key corresponding to the scalar value s of a
// private key.
func (s *PrivKeyScalar) Public() *PublicKey {
	p := NewPoint().Mul((*big.Int)(s), B8)
	pk := PublicKey(*p)
	return &pk
}

// BigInt returns the big.Int corresponding to a PrivKeyScalar.
func (s *PrivKeyScalar) BigInt() *big.Int {
	return (*big.Int)(s)
}

// ****
// * Public Key
// ****

// PublicKey represents an EdDSA public key, which is a curve point.
type PublicKey Point

// MarshalText implements the marshaler for PublicKey
func (pk PublicKey) MarshalText() ([]byte, error) {
	pkc := pk.Compress()
	return utils.Hex(pkc[:]).MarshalText()
}

// String returns the string representation of the PublicKey
func (pk PublicKey) String() string {
	pkc := pk.Compress()
	return utils.Hex(pkc[:]).String()
}

// UnmarshalText implements the unmarshaler for the PublicKey
func (pk *PublicKey) UnmarshalText(h []byte) error {
	var pkc PublicKeyComp
	if err := utils.HexDecodeInto(pkc[:], h); err != nil {
		return err
	}
	pkd, err := pkc.Decompress()
	if err != nil {
		return err
	}
	*pk = *pkd
	return nil
}

// Point returns the Point corresponding to a PublicKey.
func (pk *PublicKey) Point() *Point {
	return (*Point)(pk)
}

// ****
// * Compress Public Key
// ****

// PublicKeyComp represents a compressed EdDSA Public key; it's a compressed curve
// point.
type PublicKeyComp [32]byte

// MarshalText implements the marshaler for the PublicKeyComp
func (pkComp PublicKeyComp) MarshalText() ([]byte, error) {
	return utils.Hex(pkComp[:]).MarshalText()
}

// String returns the string representation of the PublicKeyComp
func (pkComp PublicKeyComp) String() string { return utils.Hex(pkComp[:]).String() }

// UnmarshalText implements the unmarshaler for the PublicKeyComp
func (pkComp *PublicKeyComp) UnmarshalText(h []byte) error {
	return utils.HexDecodeInto(pkComp[:], h)
}

// Compress returns the PublicKeyComp for the given PublicKey
func (pk *PublicKey) Compress() PublicKeyComp {
	return PublicKeyComp((*Point)(pk).Compress())
}

// Decompress returns the PublicKey for the given PublicKeyComp
func (pkComp *PublicKeyComp) Decompress() (*PublicKey, error) {
	point, err := NewPoint().Decompress(*pkComp)
	if err != nil {
		return nil, err
	}
	pk := PublicKey(*point)
	return &pk, nil
}

// *******
// * Signature
// *******

// Signature represents an EdDSA uncompressed signature.
type Signature struct {
	R8 *Point
	S  *big.Int
}

// SignatureComp represents a compressed EdDSA signature.
type SignatureComp [64]byte

// MarshalText implements the marshaller for the SignatureComp
func (sComp SignatureComp) MarshalText() ([]byte, error) {
	return utils.Hex(sComp[:]).MarshalText()
}

// String returns the string representation of the SignatureComp
func (sComp SignatureComp) String() string { return utils.Hex(sComp[:]).String() }

// UnmarshalText implements the unmarshaler for the SignatureComp
func (sComp *SignatureComp) UnmarshalText(h []byte) error {
	return utils.HexDecodeInto(sComp[:], h)
}

// Compress an EdDSA signature by concatenating the compression of
// the point R8 and the Little-Endian encoding of S.
func (s *Signature) Compress() SignatureComp {
	R8p := s.R8.Compress()
	Sp := utils.BigIntLEBytes(s.S)
	buf := [64]byte{}
	copy(buf[:32], R8p[:])
	copy(buf[32:], Sp[:])
	return SignatureComp(buf)
}

// Decompress a compressed signature into s, and also returns the decompressed
// signature.  Returns error if the Point decompression fails.
func (s *Signature) Decompress(buf [64]byte) (*Signature, error) {
	R8p := [32]byte{}
	copy(R8p[:], buf[:32])
	var err error
	if s.R8, err = NewPoint().Decompress(R8p); err != nil {
		return nil, err
	}
	s.S = utils.SetBigIntFromLEBytes(new(big.Int), buf[32:])
	return s, nil
}

// Decompress a compressed signature.  Returns error if the Point decompression
// fails.
func (sComp *SignatureComp) Decompress() (*Signature, error) {
	return new(Signature).Decompress(*sComp)
}

// SignPoseidon signs a message encoded as a big.Int in Zq
func (k *PrivateKey) SignPoseidon(msg *big.Int) *Signature {
	A := k.Public().Point() // A = kG

	kBuf := utils.BigIntLEBytes(k.Scalar().BigInt())
	h1 := Blake512(kBuf[:])

	msgBuf := utils.BigIntLEBytes(msg)
	msgBuf32 := [32]byte{}
	copy(msgBuf32[:], msgBuf[:])

	rBuf := Blake512(append(h1[32:], msgBuf32[:]...))
	r := utils.SetBigIntFromLEBytes(new(big.Int), rBuf) // r = H(H_{32..63}(k), msg)
	r.Mod(r, SubOrder)
	R8 := NewPoint().Mul(r, B8) // R8 = r * B8

	hmInput := []*big.Int{R8.X, R8.Y, A.X, A.Y, msg}
	poseidonParam := poseidon.NewParams(6, 6, 52)
	hm := poseidon.Hash(hmInput, poseidonParam) // hm = H1(8*R.x, 8*R.y, A.x, A.y, msg)

	S := k.Scalar().BigInt()
	S = S.Mul(hm, S)
	S.Add(r, S)
	S.Mod(S, SubOrder) // S = r + hm * k

	return &Signature{R8: R8, S: S}
}

// VerifyPoseidon verifies the signature of a message encoded as a big.Int in Fq
func (pk *PublicKey) VerifyPoseidon(msg *big.Int, sig *Signature) bool {
	hmInput := []*big.Int{sig.R8.X, sig.R8.Y, pk.X, pk.Y, msg}
	poseidonParam := poseidon.NewParams(6, 6, 52)

	hm := poseidon.Hash(hmInput, poseidonParam) // hm = H1(R8.x, R8.y, A.x, A.y, msg)

	left := NewPoint().Mul(sig.S, B8)       // left = s * B8
	right := NewPoint().Mul(hm, pk.Point()) // hm *
	rightProj := right.Projective()
	rightProj.Add(sig.R8.Projective(), rightProj) // right = R8 + hm * A
	right = rightProj.Affine()

	return (left.X.Cmp(right.X) == 0) && (left.Y.Cmp(right.Y) == 0)
}
