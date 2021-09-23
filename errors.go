package ethclient

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNoAnyKeyStores       = errors.New("No any keystores")
	ErrMessagePrivateKeyNil = errors.New("PrivateKey is nil")
)

type EVMErr struct {
	TxHash common.Hash // Empty if do call message.
	Err    string
}

func (e EVMErr) Error() string {
	return fmt.Sprintf("tx %v reverted reason: %v", e.TxHash.Hex(), e.Err)
}
