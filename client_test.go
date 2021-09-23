package ethclient

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/TheStarBoys/ethclient/contracts"
	"github.com/TheStarBoys/ethtypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

var (
	privateKey, _ = crypto.HexToECDSA("9a01f5c57e377e0239e6036b7b2d700454b760b2dab51390f1eeb2f64fe98b68")
	addr          = crypto.PubkeyToAddress(privateKey.PublicKey)
)

func deployTestContract(t *testing.T, ctx context.Context, client *Client) (common.Address, *types.Transaction, *contracts.Contracts, error) {
	auth, err := client.MessageToTransactOpts(ctx, Message{PrivateKey: privateKey})
	if err != nil {
		t.Fatal(err)
	}

	return contracts.DeployContracts(auth, client.RawClient())
}

func newTestClient(t *testing.T) *Client {
	backend, _ := NewTestEthBackend(privateKey, core.GenesisAlloc{
		addr: core.GenesisAccount{
			Balance: new(big.Int).Mul(big.NewInt(1000), ethtypes.Kether),
		},
	})
	// defer backend.Close()

	rpcClient, _ := backend.Attach()
	client, err := NewClient(rpcClient)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func TestBatchSendMsg(t *testing.T) {
	log.Root().SetHandler(log.DiscardHandler())
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	mesgs := make(chan Message)
	txs, errs := client.BatchSendMsg(ctx, mesgs)
	go func() {
		for i := 0; i < 5; i++ {
			to := common.HexToAddress("0x06514D014e997bcd4A9381bF0C4Dc21bD32718D4")
			mesgs <- Message{
				PrivateKey: privateKey,
				To:         &to,
			}
			t.Log("Write MSG to channel")
		}

		t.Log("Close send channel")
		close(mesgs)
	}()

	for tx := range txs {
		js, _ := tx.MarshalJSON()
		err := <-errs
		log.Info("Get Transaction", "tx", string(js), "err", err)
		assert.Equal(t, nil, err)
	}
	t.Log("Exit")
}

func TestCallContract(t *testing.T) {
	log.Root().SetHandler(log.DiscardHandler())
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Deploy Test contract.
	contractAddr, txOfContractCreation, contract, err := deployTestContract(t, ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("TestContract creation transaction", "txHex", txOfContractCreation.Hash().Hex(), "contract", contractAddr.Hex())

	contains, err := client.ConfirmTx(txOfContractCreation.Hash(), 2, 5*time.Second)
	if err != nil {
		t.Fatalf("Deploy Contract err: %v", err)
	}
	assert.Equal(t, true, contains)

	// Call contract method `testFunc1` id -> 0x88655d98
	contractAbi := contracts.GetTestContractABI()

	arg1 := "hello"
	arg2 := big.NewInt(100)
	arg3 := []byte("world")
	data, err := client.NewMethodData(contractAbi, "testFunc1", arg1, arg2, arg3)
	if err != nil {
		t.Fatal(err)
	}

	if code, err := client.RawClient().CodeAt(ctx, contractAddr, nil); err != nil || len(code) == 0 {
		t.Fatal("no code or has err: ", err)
	}

	// contract.TestFunc1(nil)
	_, err = client.CallMsg(ctx, Message{
		From: crypto.PubkeyToAddress(privateKey.PublicKey),
		To:   &contractAddr,
		Data: data,
	}, nil)
	if err != nil {
		t.Fatalf("CallMsg err: %v", err)
	}

	contractCallTx, err := client.SendMsg(ctx, Message{
		PrivateKey: privateKey,
		To:         &contractAddr,
		Data:       data,
	})
	if err != nil {
		t.Fatalf("Send single Message err: %v", err)
	}

	t.Log("contractCallTx send sucessul", "txHash", contractCallTx.Hash().Hex())

	contains, err = client.ConfirmTx(contractCallTx.Hash(), 2, 20*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, contains)

	receipt, err := client.RawClient().TransactionReceipt(ctx, contractCallTx.Hash())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Receipt", "status", receipt.Status)

	counter, err := contract.Counter(nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint64(1), counter.Uint64())
}

func TestContractRevert(t *testing.T) {
	log.Root().SetHandler(log.DiscardHandler())
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Deploy Test contract.
	contractAddr, txOfContractCreation, _, err := deployTestContract(t, ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("TestContract creation transaction", "txHex", txOfContractCreation.Hash().Hex(), "contract", contractAddr.Hex())

	contains, err := client.ConfirmTx(txOfContractCreation.Hash(), 2, 5*time.Second)
	if err != nil {
		t.Fatalf("Deploy Contract err: %v", err)
	}
	assert.Equal(t, true, contains)

	// Call contract method `testFunc1` id -> 0x88655d98
	contractAbi := contracts.GetTestContractABI()
	data, err := client.NewMethodData(contractAbi, "testReverted")
	assert.Equal(t, nil, err)

	// Send successful, but executation failed.
	contractCallTx, err := client.SendMsg(ctx, Message{
		PrivateKey: privateKey,
		To:         &contractAddr,
		Data:       data,
		Gas:        210000,
		GasPrice:   big.NewInt(10),
	})
	if err != nil {
		t.Fatalf("Send single Message, err: %v", err)
	}

	client.ConfirmTx(contractCallTx.Hash(), 1, 2*time.Second)
	receipt, err := client.rawClient.TransactionReceipt(ctx, contractCallTx.Hash())
	assert.Equal(t, nil, err)
	assert.Equal(t, types.ReceiptStatusFailed, receipt.Status)

	t.Log("contractCallTx send sucessul", "txHash", contractCallTx.Hash().Hex())

	// Send failed, because estimateGas faield.
	contractCallTx, err = client.SendMsg(ctx, Message{
		PrivateKey: privateKey,
		To:         &contractAddr,
		Data:       data,
	})
	t.Log("Send Message without specific gas and gasPrice, err: ", err)
	// Send Message without specific gas and gasPrice, err:  NewTransaction err: execution reverted: test reverted
	assert.NotEqual(t, nil, err, "expect revert transaction")

	// Call failed, because evm execution faield.
	returnData, err := client.CallMsg(ctx, Message{
		PrivateKey: privateKey,
		To:         &contractAddr,
		Data:       data,
	}, nil)
	t.Log("Call Message err: ", err)
	assert.Equal(t, 0, len(returnData))
	assert.NotEqual(t, nil, err, "expect revert transaction")
}
