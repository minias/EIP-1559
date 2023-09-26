package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/minias/EIP-1559/env"
	"github.com/minias/EIP-1559/eth"
)

func init() {
	env.ReadConfig(env.InitProfile())
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("useage: go run main.go [privateKeyHex] [toAddressHex] [AmountWei]")
		fmt.Printf("example: go run main.go %s %s %s\n",
			"0x12345678901234567890123456789012",
			"0xd2716D0d298284Dc955090A03ba16a916B219fA6",
			"1100000000000")
		return
	}

	// Data binding
	privateKeyHex := os.Args[1]
	recipientStr := os.Args[2]
	value := big.NewInt(0)
	value.SetString(os.Args[3], 10) // 100 * 1e18

	// Eth Transaction Transfer Case
	txHash, err := eth.SendEtherEIP1559(privateKeyHex, recipientStr, nil)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	fmt.Printf("Transfer Scan: %s%s\n", env.Conf.BlockChain.SCANCURL, txHash)

	// ERC-20 Transaction Transfer Case
	data := eth.SetTokenData(common.HexToAddress(recipientStr), value)
	fmt.Printf("data: %x\n", data)
	txHash, err = eth.SendEtherEIP1559(privateKeyHex, "", data)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	fmt.Printf("Transfer Scan: %s%s\n", env.Conf.BlockChain.SCANCURL, txHash)
}
