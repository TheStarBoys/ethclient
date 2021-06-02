package ethclient

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

type Client struct {
	RawClient *ethclient.Client
	Subscriber
}

func Dial(rawurl string) (*Client, error) {
	c, err := ethclient.Dial(rawurl)
	if err != nil {
		return nil, err
	}

	subscriber, err := NewChainSubscriber(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		RawClient:  c,
		Subscriber: subscriber,
	}, nil
}

func NewClient(c *ethclient.Client) (*Client, error) {
	subscriber, err := NewChainSubscriber(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		RawClient:  c,
		Subscriber: subscriber,
	}, nil
}

type Message struct {
	From       common.Address    // the sender of the 'transaction'
	PrivateKey *ecdsa.PrivateKey // overwrite From if not nil
	To         *common.Address   // the destination contract (nil for contract creation)
	Gas        uint64            // if 0, the call executes with near-infinite gas
	GasPrice   *big.Int          // wei <-> gas exchange ratio
	Value      *big.Int          // amount of wei sent along with the call
	Data       []byte            // input data, usually an ABI-encoded contract method invocation

	AccessList types.AccessList // EIP-2930 access list.
}

func (c *Client) NewMethodData(a abi.ABI, methodName string, args ...interface{}) ([]byte, error) {
	return a.Pack(methodName, args...)
}

func (c *Client) BatchSendMsg(ctx context.Context, msgs <-chan Message) (<-chan *types.Transaction, <-chan error) {
	txs := make(chan *types.Transaction, 10)
	errs := make(chan error, 10)
	go func() {
		for msg := range msgs {
			log.Info("Receive Msg from channel", "msg", msg)
			tx, err := c.SendMsg(ctx, msg)
			txs <- tx
			errs <- err
			log.Info("Send Msg successful", "tx", tx.Hash().Hex())
		}

		log.Info("Close BatchSendMsg channel")
		close(txs)
		close(errs)
	}()
	return txs, errs
}

func (c *Client) SendMsg(ctx context.Context, msg Message) (*types.Transaction, error) {
	if msg.PrivateKey == nil {
		return nil, fmt.Errorf("PrivateKey is nil")
	}

	msg.From = crypto.PubkeyToAddress(msg.PrivateKey.PublicKey)

	ethMesg := ethereum.CallMsg{
		From:       msg.From,
		To:         msg.To,
		Gas:        msg.Gas,
		GasPrice:   msg.GasPrice,
		Value:      msg.Value,
		Data:       msg.Data,
		AccessList: msg.AccessList,
	}

	tx, err := c.NewTransaction(ctx, ethMesg)
	if err != nil {
		return nil, err
	}

	chainID, err := c.RawClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP2930Signer(chainID), msg.PrivateKey)
	if err != nil {
		return nil, err
	}

	err = c.RawClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (c *Client) NewTransaction(ctx context.Context, msg ethereum.CallMsg) (*types.Transaction, error) {
	if msg.Gas == 0 {
		gas, err := c.RawClient.EstimateGas(ctx, msg)
		if err != nil {
			return nil, err
		}

		msg.Gas = gas
	}

	if msg.GasPrice == nil || msg.GasPrice.Uint64() == 0 {
		var err error
		msg.GasPrice, err = c.RawClient.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
	}

	// TODO: nonce manager
	nonce, err := c.RawClient.PendingNonceAt(ctx, msg.From)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, *msg.To, msg.Value, msg.Gas, msg.GasPrice, msg.Data)

	return tx, nil
}

func (c *Client) ConfirmTx(txHash common.Hash, n uint, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// use SubscribeNewHead to confirm the signed transaction was contained in the new block.
	headerChan := make(chan *types.Header)
	err := c.SubscribeNewHead(ctx, headerChan)
	if err != nil {
		return false, err
	}

	for {
		select {
		case header := <-headerChan:
			block, err := c.RawClient.BlockByHash(ctx, header.Hash())
			if err != nil {
				return false, err
			}

			gotTx := block.Transaction(txHash)
			if gotTx != nil {
				// got tx receipt
				// TODO: implement wait n blocks
				return true, nil
			}
		case <-ctx.Done():
			// Not in chain
			return false, nil
		}
	}
}

func (c *Client) MessageToTransactOpts(ctx context.Context, msg Message) (*bind.TransactOpts, error) {
	if msg.PrivateKey == nil {
		return nil, fmt.Errorf("PrivateKey is nil")
	}
	msg.From = crypto.PubkeyToAddress(msg.PrivateKey.PublicKey)

	nonce, err := c.RawClient.PendingNonceAt(ctx, msg.From)
	if err != nil {
		return nil, err
	}

	signFunc := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		chainID, err := c.RawClient.ChainID(ctx)
		if err != nil {
			return nil, err
		}

		return types.SignTx(tx, types.NewEIP2930Signer(chainID), msg.PrivateKey)
	}

	auth := bind.NewKeyedTransactor(msg.PrivateKey)
	auth.Signer = signFunc
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = msg.Value  // in wei
	auth.GasLimit = msg.Gas // in units
	auth.GasPrice = msg.GasPrice

	return auth, nil
}
