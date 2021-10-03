package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-panta/bank-smart-contract/internal/contract/bank"
)

func main() {
	log.SetFlags(log.Lshortfile)

	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("a4ff3a5a34cc78807edfad12c41e1cf965a89f81e125b1386672e6190e2bfaa8")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.PublicKey
	log.Printf("public_key: %+v", publicKey)

	fromAddress := crypto.PubkeyToAddress(publicKey)
	log.Printf("from_address: %v", fromAddress)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("nonce: %d", nonce)

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("gas_price: %d", gasPrice)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	callMsg := ethereum.CallMsg{
		From: common.HexToAddress(fromAddress.Hex()),
		Data: common.FromHex(bank.BankBin),
	}
	gasLimit, err := client.EstimateGas(ctx, callMsg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("gas_limit: %d", gasLimit)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	txAddress, tx, instance, err := bank.DeployBank(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tx_address: %s", txAddress.Hex())
	log.Printf("tx: %s", tx.Hash().Hex())

	nonce++
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(1000000000000000000) // in wei
	depositTx, err := instance.Deposit(auth)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("deposit_tx: %s", depositTx.Hash().Hex())

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("block_number: %d", blockNumber)

	callOpts := &bind.CallOpts{
		Pending:     true,
		From:        fromAddress,
		BlockNumber: big.NewInt(int64(blockNumber)),
		Context:     ctx,
	}
	depositBalance, err := instance.GetBalance(callOpts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("deposit_balance: %d", depositBalance)

	bankTVL, err := client.PendingBalanceAt(ctx, txAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bank_tvl: %d", bankTVL)

	nonce++
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	withdrawTx, err := instance.Withdraw(auth, fromAddress, big.NewInt(100000000000000000))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("withdraw_tx: %s", withdrawTx.Hash().Hex())

	bankTVL, err = client.PendingBalanceAt(ctx, txAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bank_tvl: %d", bankTVL)
}
