package eth_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/minias/EIP-1559/eth"
)

// TestSendEtherEIP1559 godoc
func TestSendEtherEIP1559(t *testing.T) {
	privateKeyHex := ""
	toStr := ""
	value := big.NewInt(0)
	value.SetString("", 10)
	data := eth.SetTokenData(common.HexToAddress(toStr), value)
	s, err := eth.SendEtherEIP1559(privateKeyHex, toStr, data)
	if err != nil {
		t.Error("Wrong result")
	}
	if len(s) == 66 {
		t.Log("running TestAdd")
	}
}

// TESTSendNew1 godoc
func TESTSendNew1(t *testing.T) {
	privateKeyHex := ""
	toStr := ""
	s, err := eth.SendNew1(privateKeyHex, toStr)
	if err != nil {
		t.Error("Wrong result")
	}
	if len(s) == 66 {
		t.Log("running TestAdd")
	}
}

// TESTSendNew2 godoc
func TESTSendNew2(t *testing.T) {
	ContractStr := ""
	privateKeyHex := ""
	toStr := ""
	s, err := eth.SendNew2(ContractStr, privateKeyHex, toStr)
	if err != nil {
		t.Error("Wrong result")
	}
	if len(s) == 66 {
		t.Log("running TestAdd")
	}
}

// TESTSendNew3 godoc
func TESTSendNew3(t *testing.T) {
	privateKeyHex := ""
	toStr := ""
	s, err := eth.SendNew3(privateKeyHex, toStr)
	if err != nil {
		t.Error("Wrong result")
	}
	if len(s) == 66 {
		t.Log("running TestAdd")
	}
}

// TESTSendNew4 godoc
func TESTSendNew4(t *testing.T) {
	ContractStr := ""
	privateKeyHex := ""
	toStr := ""
	s, err := eth.SendNew4(ContractStr, privateKeyHex, toStr)
	if err != nil {
		t.Error("Wrong result")
	}
	if len(s) == 66 {
		t.Log("running TestAdd")
	}
}
