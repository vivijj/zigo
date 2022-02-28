package state

type OpError struct {
	msg string
}

func (e *OpError) Error() string {
	return e.msg
}

// InvalidToken deposit error
var InvalidToken = &OpError{"deposit token is out of range, this should be enforced by contract"}

var (
	InvalidTokenId      = &OpError{"token id is not supported"}
	InvalidFeeTokenId   = &OpError{"fee token id is not supported"}
	FromAccountNotFound = &OpError{"from account does not exist"}
	FromAccountLocked   = &OpError{"account is locked"}
	NonceMismatch       = &OpError{"nonce mismatch"}
	InsufficientBalance = &OpError{"not enough balance"}
	InvalidSignature    = &OpError{"L2 signature is incorrect"}
)

var (
	TargetAccountZero        = &OpError{"transfer to account with address 0 is not allowed"}
	TransferAccountIncorrect = &OpError{"transfer account id is incorrect"}
)

var (
	FromAccountIncorrect = &OpError{"withdraw account id is incorrect"}
)

// pubkey update

var (
	AccountNotFound       = &OpError{"account does not exist"}
	InvalidAuthData       = &OpError{"pubkeyUpdate L1 auth data is incorrect"}
	InvalidAccountAddress = &OpError{"account address is incorrect"}
	InvalidAccountId      = &OpError{"account id is incorrect"}
)
