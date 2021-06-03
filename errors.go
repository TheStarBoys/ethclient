package ethclient

import "errors"

var (
	ErrNoAnyKeyStores       = errors.New("No any keystores")
	ErrMessagePrivateKeyNil = errors.New("PrivateKey is nil")
)
