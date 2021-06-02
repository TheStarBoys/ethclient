package ethclient

import (
	"bytes"
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/TheStarBoys/ethclient/contracts"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	log.Root().SetHandler(log.StdoutHandler)
	log.Info("Dial.....")
	client, err := Dial("ws://localhost:8546")
	if err != nil {
		t.Fatal(err)
	}

	log.Info("Dial successful!")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	privateKey, _ := crypto.HexToECDSA("9a01f5c57e377e0239e6036b7b2d700454b760b2dab51390f1eeb2f64fe98b68")

	// Deploy Test contract.
	auth, err := client.MessageToTransactOpts(ctx, Message{PrivateKey: privateKey})
	if err != nil {
		t.Fatal(err)
	}

	contractAddr, txOfContractCreation, contract, err := contracts.DeployContracts(auth, client.RawClient)
	if err != nil {
		t.Fatal(err)
	}

	client.ConfirmTx(txOfContractCreation.Hash(), 6, 20*time.Second)

	log.Info("TestContract creation transaction", "txHex", txOfContractCreation.Hash().Hex(), "contract", contractAddr.Hex())

	// Call contract method `testFunc1` id -> 0x88655d98
	contractAbi, err := abi.JSON(bytes.NewBuffer([]byte(contracts.ContractsABI)))
	if err != nil {
		t.Fatal(err)
	}

	methodId := common.FromHex("0x88655d98")
	arg1 := "hello"
	arg2 := big.NewInt(100)
	arg3 := []byte("world")
	data, err := client.NewMethodData(contractAbi, "testFunc1", arg1, arg2, arg3)
	if err != nil {
		t.Fatal(err)
	}

	log.Info("Method data", "data", common.Bytes2Hex(data))

	if code, err := client.RawClient.CodeAt(ctx, contractAddr, nil); err != nil || len(code) == 0 {
		t.Fatal("no code or has err: ", err)
	}

	// contract.TestFunc1(nil)
	_, err = client.RawClient.CallContract(ctx, ethereum.CallMsg{
		From: crypto.PubkeyToAddress(privateKey.PublicKey),
		To:   &contractAddr,
		Data: data,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	contractCallTx, err := client.SendMsg(ctx, Message{
		PrivateKey: privateKey,
		To:         &contractAddr,
		Data:       data,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info("contractCallTx send sucessul", "methodId", common.Bytes2Hex(methodId), "txHash", contractCallTx.Hash().Hex())

	contains, err := client.ConfirmTx(contractCallTx.Hash(), 6, 20*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := client.RawClient.TransactionReceipt(ctx, contractCallTx.Hash())
	if err != nil {
		t.Fatal(err)
	}
	log.Info("Receipt", "status", receipt.Status)

	log.Info("weather transaction contains in chain", "contains", contains)

	counter, err := contract.Counter(nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint64(1), counter.Uint64())

	mesgs := make(chan Message)
	txs, errs := client.BatchSendMsg(ctx, mesgs)
	go func() {
		for i := 0; i < 5; i++ {
			to := common.HexToAddress("0x06514D014e997bcd4A9381bF0C4Dc21bD32718D4")
			mesgs <- Message{
				PrivateKey: privateKey,
				To:         &to,
			}
			log.Info("Write MSG to channel")
		}

		log.Info("Close send channel")
		close(mesgs)
	}()

	for tx := range txs {
		js, _ := tx.MarshalJSON()
		log.Info("Get Transaction", "tx", string(js), "err", <-errs)
	}
	log.Info("Exit")
}
