package contracts

import (
	"bytes"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func GetTestContractABI() abi.ABI {
	contractAbi, err := abi.JSON(bytes.NewBuffer([]byte(ContractsABI)))
	if err != nil {
		panic(err)
	}

	return contractAbi
}
