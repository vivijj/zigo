package transaction

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/vivijj/zigo/crypto/babyjub"
)

// PubkeyUpdateTx will set the owner's public key associated with the account.
// without public key set, account is unable to execute any L2 transactions.
type PubkeyUpdateTx struct {
	AccountId  int
	Nonce      int
	ValidUntil int
	FeeToken   int
	Fee        *big.Int
	Account    common.Address
	PubKey     babyjub.PublicKey
	AuthData   []byte
}

func (tx PubkeyUpdateTx) isZionTx() {}

func (tx PubkeyUpdateTx) GetBytes() (out []byte) {

	out = append(out, []byte(PubKeyUpdate)...)
	out = append(out, IntToBytes(tx.AccountId)...)
	out = append(out, IntToBytes(tx.Nonce)...)
	out = append(out, IntToBytes(tx.ValidUntil)...)
	out = append(out, IntToBytes(tx.FeeToken)...)
	out = append(out, tx.Fee.Bytes()...)
	out = append(out, tx.Account.Bytes()...)
	out = append(out, tx.PubKey.X.Bytes()...)
	out = append(out, tx.PubKey.Y.Bytes()...)

	return
}

func (tx PubkeyUpdateTx) IsAuthDataValid() bool {
	userAddr := tx.Account
	sig := tx.AuthData
	msgHash := tx.hashTypedData()
	fmt.Println("msgHash is: ", msgHash)
	if len(sig) != 65 {
		return false
	}
	if sig[64] != 27 && sig[64] != 28 {
		return false
	}
	sig[64] -= 27
	pubkey, err := crypto.SigToPub(msgHash, sig)
	if err != nil {
		return false

	}
	recoverAddr := crypto.PubkeyToAddress(*pubkey)
	fmt.Println("recover address is: ", recoverAddr)
	if !bytes.Equal(recoverAddr.Bytes(), userAddr.Bytes()) {
		return false
	}
	return true
}

func (tx PubkeyUpdateTx) hashTypedData() []byte {
	typedData := tx.TypedData()
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		panic(err)
		fmt.Println("fail to hash domain: ", err)
	}
	fmt.Println("domain is :", domainSeparator)
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		fmt.Println("fail to hash typed data: ", err)
	}
	fmt.Println("typed data is ", typedDataHash)
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return crypto.Keccak256(rawData)
}

// TODO(vivijj): uint96 is not a valid type in the hashstruct
func (tx PubkeyUpdateTx) TypedData() apitypes.TypedData {
	return apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PubkeyUpdate": []apitypes.Type{
				{"address", "address"},
				{"accountId", "uint32"},
				{"feeTokenId", "uint16"},
				{"fee", "uint96"},
				{"publicKey", "uint256"},
				{"validUntil", "uint32"},
				{"nonce", "uint32"},
			},
		},
		PrimaryType: "PubkeyUpdate",
		Domain: apitypes.TypedDataDomain{
			Name:              "Zion Rollup",
			Version:           "1.0.0",
			ChainId:           (*math.HexOrDecimal256)(big.NewInt(43)),
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Message: apitypes.TypedDataMessage{
			"address":    tx.Account.String(),
			"accountId":  (*math.HexOrDecimal256)(big.NewInt(int64(tx.AccountId))),
			"feeTokenId": (*math.HexOrDecimal256)(big.NewInt(int64(tx.FeeToken))),
			"fee":        (*math.HexOrDecimal256)(tx.Fee),
			"publicKey":  (*math.HexOrDecimal256)((&(tx.PubKey)).CompressBi()),
			"validUntil": (*math.HexOrDecimal256)(big.NewInt(int64(tx.ValidUntil))),
			"nonce":      (*math.HexOrDecimal256)(big.NewInt(int64(tx.Nonce))),
		},
	}
}
