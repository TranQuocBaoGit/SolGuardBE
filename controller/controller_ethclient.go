package controller

import (
	"context"
	"encoding/hex"
	"fmt"
	"getContractDeployment/helper"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getBlockByNumber(ctx context.Context, client *ethclient.Client, num *big.Int) *types.Block{
    if (num.Int64() == -1){
        block, err := client.BlockByNumber(context.Background(), nil)
        helper.CheckError(err)
        return block
    }
    block, err := client.BlockByNumber(context.Background(), num)
    helper.CheckError(err)
    return block
}

func GetContractAddressFromToBlock(ctx context.Context, client *ethclient.Client, num *big.Int) []common.Address{
    // Get block
    block := getBlockByNumber(ctx, client, num)

    allTx := []common.Hash{}
    currentBlockNumber := block.Number().Int64()
    fmt.Println("Current block number: ", currentBlockNumber)

    for i := 20; i >= 0; i--{
        block := getBlockByNumber(ctx, client, big.NewInt(currentBlockNumber - int64(i)))
        fmt.Println("Get block ", block.Number().String())
        fmt.Println("Contract creation transaction in this block: ")
        for _, tx := range block.Transactions(){
            if tx.To() == nil{
                fmt.Println(tx.Hash())
                allTx = append(allTx, tx.Hash())
            }
        }
        fmt.Println("")
        time.Sleep(time.Second * 3)
    }

    fmt.Println("All contract deploy transaction: ", allTx)

    time.Sleep(time.Second * 5)

    allContractAddress := []common.Address{}
    for _, tx := range allTx{
        receipt, err := client.TransactionReceipt(context.Background(), tx)
        helper.CheckError(err)
        allContractAddress = append(allContractAddress, receipt.ContractAddress)
        time.Sleep(time.Second * 1)
    }
    return allContractAddress
}


func getSmartContractByteCode(client *ethclient.Client, address string) string{
    contractAddress := common.HexToAddress(address)
    bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
    helper.CheckError(err)

    bytecodeStr := hex.EncodeToString(bytecode)

    return bytecodeStr
}