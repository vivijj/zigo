package transaction

import (
	"encoding/json"
	"fmt"
)

// a simple helper to marshal and unmarshal ZionTx interface

type ZionTxJson struct {
	Type  TransactionType
	Value json.RawMessage
}

func FromZionTxToJson(txData ZionTx) ZionTxJson {
	var txType TransactionType
	switch txData.(type) {
	case TransferTx:
		txType = Transfer
	case WithdrawTx:
		txType = Withdraw
	case PubkeyUpdateTx:
		txType = PubKeyUpdate
	}
	jtx, _ := json.Marshal(txData)
	return ZionTxJson{
		Type:  txType,
		Value: jtx,
	}
}

// ParseZionTx Parse Zion tx from the ZionTxJson, if failed, return the nil interface.
func (jtx *ZionTxJson) ParseZionTx() ZionTx {
	switch jtx.Type {
	case Withdraw:
		var withdraw WithdrawTx
		_ = json.Unmarshal(jtx.Value, &withdraw)
		return withdraw

	case Transfer:
		var transfer TransferTx
		_ = json.Unmarshal(jtx.Value, &transfer)
		return transfer

	case PubKeyUpdate:
		var pubkeyUpdate PubkeyUpdateTx
		_ = json.Unmarshal(jtx.Value, &pubkeyUpdate)
		return pubkeyUpdate
	}
	return nil
}

type TypedZionTx struct {
	Type  TransactionType
	Value ZionTx
}

func (tx *TypedZionTx) UnmarshalJSON(data []byte) error {
	var txJSon ZionTxJson
	err := json.Unmarshal(data, &txJSon)
	switch txJSon.Type {
	case Withdraw:
		var withdraw WithdrawTx
		err = json.Unmarshal(txJSon.Value, &withdraw)
		tx.Value = withdraw
	case Transfer:
		var transfer TransferTx
		err = json.Unmarshal(txJSon.Value, &transfer)
		tx.Value = transfer
	case PubKeyUpdate:
		var pubkeyUpdate PubkeyUpdateTx
		err = json.Unmarshal(txJSon.Value, &pubkeyUpdate)
		tx.Value = pubkeyUpdate
	default:
		return fmt.Errorf("invalid ZionTx type")

	}
	return err
}
