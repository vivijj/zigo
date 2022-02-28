package witness

import (
	"github.com/vivijj/zigo/pkg/operation"
	"github.com/vivijj/zigo/types/fr"
)

// This serde include some data that will be serialized as json to prover.

// ProverData is data prover needs to calculate proof of the given block.
type ProverData struct {
	MerkleRootBefore      fr.Fr          `json:"merkle_root_before"`
	MerkleRootAfter       fr.Fr          `json:"merkle_root_after"`
	TimeStamp             fr.Fr          `json:"time_stamp"`
	OperatorId            fr.Fr          `json:"operator_id"`
	AccountUpdateOperator AccountWitness `json:"account_update_operator"`
	TxWitness             []Operation    `json:"tx_witness"`
}
type Operation struct {
	TxType  string
	Tx      operation.ZionOp
	Witness Witness
}

type Witness struct {
	SignatureFrom     *SigValue `json:"signature_from"`
	AccountMerkleRoot fr.Fr     `json:"account_merkle_root"`

	BalanceUpdateFrom    BalanceWitness `json:"balance_update_from"`
	BalanceUpdateFeeFrom BalanceWitness `json:"balance_update_fee_from"`
	AccountUpdateFrom    AccountWitness `json:"account_update_from"`

	BalanceUpdateTo BalanceWitness `json:"balance_update_to"`
	AccountUpdateTo AccountWitness `json:"account_update_to"`

	BalanceUpdateOperator BalanceWitness `json:"balance_update_operator"`
	AccountUpdateOperator AccountWitness `json:"account_update_operator"`

	NumConditionalTransactionAfter fr.Fr `json:"num_conditional_transaction_after"`
}

type AccountWitness struct {
	AccountId     fr.Fr        `json:"account_id"`
	AccountBefore AccountValue `json:"account_before"`
	AccountAfter  AccountValue `json:"account_after"`
	RootBefore    fr.Fr        `json:"root_before"`
	RootAfter     fr.Fr        `json:"root_after"`
	AccountPath   []fr.Fr      `json:"proof"`
}

type BalanceWitness struct {
	TokenId       fr.Fr        `json:"token_id"`
	BalanceBefore BalanceValue `json:"balance_before"`
	BalanceAfter  BalanceValue `json:"balance_after"`
	RootBefore    fr.Fr        `json:"root_before"`
	RootAfter     fr.Fr        `json:"root_after"`
	Proof         []fr.Fr      `json:"proof"`
}

type AccountValue struct {
	Address     fr.Fr `json:"address"`
	PublicKeyX  fr.Fr `json:"public_key_x"`
	PublicKeyT  fr.Fr `json:"public_key_t"`
	Nonce       fr.Fr `json:"nonce"`
	BalanceRoot fr.Fr `json:"balance_root"`
}

type BalanceValue struct {
	Balance fr.Fr `json:"balance"`
}

type SigValue struct {
	Rx fr.Fr `json:"Rx"`
	Ry fr.Fr `json:"Ry"`
	S  fr.Fr `json:"s"`
}
