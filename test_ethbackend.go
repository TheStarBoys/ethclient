package ethclient

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
)

func NewTestEthBackend(privateKey *ecdsa.PrivateKey, alloc core.GenesisAlloc) (*node.Node, error) {
	// Generate test chain.
	etherbase := crypto.PubkeyToAddress(privateKey.PublicKey)
	genesis := generateTestGenesis(etherbase, alloc)
	// Create node
	n, err := node.New(&node.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't create new node: %v", err)
	}
	// Create Ethereum Service
	config := &ethconfig.Config{Genesis: genesis}
	// config.Ethash.PowMode = ethash.ModeFake
	ethservice, err := eth.New(n, config)
	if err != nil {
		return nil, fmt.Errorf("can't create new ethereum service: %v", err)
	}
	// Import the test chain.
	if err := n.Start(); err != nil {
		return nil, fmt.Errorf("can't start test node: %v", err)
	}
	if err := saveMiner(n, privateKey); err != nil {
		return nil, fmt.Errorf("save miner err: %v", err)
	}

	ethservice.SetEtherbase(etherbase)
	err = ethservice.StartMining(1)
	if err != nil {
		return nil, fmt.Errorf("can't start mining, err: %v", err)
	}

	return n, nil
}

func saveMiner(stack *node.Node, minerPrivKey *ecdsa.PrivateKey) error {
	var ks *keystore.KeyStore
	if keystores := stack.AccountManager().Backends(keystore.KeyStoreType); len(keystores) > 0 {
		ks = keystores[0].(*keystore.KeyStore)
	} else {
		return fmt.Errorf("No any keystores")
	}

	passphrase := ""
	account, err := ks.ImportECDSA(minerPrivKey, passphrase)
	if err != nil {
		return err
	}

	return ks.Unlock(account, passphrase)
}

func generateTestGenesis(miner common.Address, alloc core.GenesisAlloc) *core.Genesis {
	// db := rawdb.NewMemoryDatabase()
	// config := params.AllEthashProtocolChanges
	genesis := core.DeveloperGenesisBlock(1, miner)
	genesis.Alloc = alloc
	// genesis := &core.Genesis{
	// 	Config:    config,
	// 	Alloc:     alloc,
	// 	ExtraData: []byte("test genesis"),
	// 	Timestamp: 9000,
	// }
	return genesis
}
