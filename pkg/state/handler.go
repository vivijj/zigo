package state

// Tx handler include the logic of each individual transaction handling.
// By transaction, we assume both zion tx and priority transaction.

func invariant(condition bool, err error) error {
	if !condition {
		return err
	}
	return nil
}
