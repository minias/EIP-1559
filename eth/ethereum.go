package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/minias/EIP-1559/env"
	"golang.org/x/crypto/sha3"
)

// Create Ethereum Client  godoc
func ethClient() (*ethclient.Client, error) {
	return ethclient.Dial(env.Conf.BlockChain.RPCURL)
}

// New Type Send Transfer EIP1559 godoc
func SendEtherEIP1559(privateKeyHex string, toStr string, data []byte) (string, error) {
	var value *big.Int
	var gasLimit uint64
	var gasTipCap *big.Int

	// Connect the eth client to the RPC nodeURL in Yml file.
	client, err := ethClient()
	if err != nil {
		return "", fmt.Errorf("ethClient: %v", err)
	}

	// Toaddress
	ToAddress := common.HexToAddress(toStr)
	// Send Transfer signer From signer Private Key String
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("HexToECDSA: %v", err)
	}
	//Get ChainId from Ethereum Network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("NetworkID: %v", err)
	}
	// Get Nonce from From address
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return "", fmt.Errorf("PendingNonceAt: %v", err)
	}
	//SuggestGasPrice //35667346 Wei
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("SuggestGasPrice: %v", err)
	}
	//gasFeeCap := GasFeeUp(gasPrice, 1.012)
	//gasFeeCap := GasFeeUp(gasPrice, 1.002)
	gasFeeCap := gasPrice
	gasTipCap = GetMaxPriorityFee(100) //1000000000

	if data != nil {
		Contract := common.HexToAddress(env.Conf.BlockChain.CONTRACT_ADDRESS)
		value = big.NewInt(0)
		//EstimateGas //34372
		EstimateGas, err := client.EstimateGas(context.Background(),
			ethereum.CallMsg{
				To:   &Contract,
				From: crypto.PubkeyToAddress(privateKey.PublicKey),
				Data: data,
			},
		)
		if err != nil {
			return "", fmt.Errorf("EstimateGas: %v", err)
		}
		fmt.Printf("calcGasFee: %v\n", CalcGasFee(gasFeeCap, EstimateGas))
		// 21000 * 1.02
		//gasLimit = GasLimitUp(EstimateGas, 1.001)
		gasLimit = EstimateGas
		gasTipCap = GetMaxPriorityFee(100) //1000000000
	} else {
		value = CalcGasFee(gasFeeCap, 21000)
		//EstimateGas // 21000 wei
		EstimateGas, err := client.EstimateGas(context.Background(),
			ethereum.CallMsg{
				To:    &ToAddress,
				From:  crypto.PubkeyToAddress(privateKey.PublicKey),
				Value: value,
			},
		)
		if err != nil {
			return "", fmt.Errorf("EstimateGas: %v", err)
		}
		// 21000 * 1.02
		gasLimit = GasLimitUp(EstimateGas, 1.001)

	}

	fmt.Printf("gasPrice: %v\n", gasPrice)   //22629304720
	fmt.Printf("GasFeeCap: %v\n", gasFeeCap) //49912709427000
	fmt.Printf("gasLimit: %v\n", gasLimit)   //300000
	fmt.Printf("value: %v\n", value)         //480917983896000 //0.000480917983896

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,   //11155111
		Nonce:     nonce,     //67
		GasTipCap: gasTipCap, //a.k.a. maxPriorityFeePerGas //1000000000
		GasFeeCap: gasFeeCap, // a.k.a. maxFeePerGas //22900856376
		Gas:       gasLimit,  //21420
		To:        &ToAddress,
		Value:     value,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("SignTx: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// CalcGasFee godoc
func CalcGasFee(bigInt *big.Int, multiplier uint64) *big.Int {
	result := new(big.Int).Set(bigInt)
	multiplierBigInt := new(big.Int).SetUint64(multiplier)
	result.Mul(result, multiplierBigInt)
	return result
}

// GasLimitUp godoc
func GasLimitUp(inputUint uint64, multiplier float64) uint64 {
	resultFloat := float64(inputUint) * multiplier
	resultUint := uint64(resultFloat)
	return resultUint
}

// GasFeeUp godoc
func GasFeeUp(bigInt *big.Int, multiplier float64) *big.Int {
	multiplierFloat := new(big.Float).SetFloat64(multiplier)
	bigIntFloat := new(big.Float).SetInt(bigInt)
	resultFloat := new(big.Float).Mul(bigIntFloat, multiplierFloat)
	resultInt := new(big.Int)
	resultFloat.Int(resultInt)
	return resultInt
}

// GetMaxPriorityFee godoc
func GetMaxPriorityFee(PriorityType int) *big.Int {
	PriorityFee := new(big.Int)
	//PriorityType
	switch PriorityType {
	case 50:
		PriorityFee.SetString("500000000", 10) //0.5Gwei
	case 100:
		PriorityFee.SetString("1000000000", 10) //1Gwei
	case 150:
		PriorityFee.SetString("1500000000", 10) //1.5Gwei
	default:
		PriorityFee.SetString("1000000000", 10) //1Gwei
	}

	fmt.Printf("PriorityFee: %v\n", PriorityFee)

	return PriorityFee
}

// SendEther godoc
func SendEther(privateKeyHex string, recipientStr string, ethAmountStr string, data []byte) (string, error) {
	client, err := ethClient()
	if err != nil {
		return "", err
	}
	// Send Transfer signer From signer Private Key String
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	//Get ChainId from Ethereum Network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	// Convert recipient address to common.Address
	recipientAddr := common.HexToAddress(recipientStr)

	// parsing the Ether amount
	ethAmount, ok := new(big.Int).SetString(ethAmountStr, 10)
	if !ok {
		return "", fmt.Errorf("ethAmountStr parsing failed: %s", ethAmountStr)
	}

	// Create Transaction
	tx, err := NewTx(crypto.PubkeyToAddress(privateKey.PublicKey), recipientAddr, ethAmount, data)
	if err != nil {
		return "", err
	}

	// Transaction Signatures
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// Transaction Transfer
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// NewTx godoc
func NewTx(from common.Address, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	// Connect the eth client to the nodeURL.
	client, err := ethClient()
	if err != nil {
		return nil, err
	}
	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return nil, err
	}
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{To: &to, From: from, Value: value, Data: data})
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	return tx, nil
}

// SetTokenData godoc
func SetTokenData(to common.Address, amount *big.Int) []byte {
	//create Signature hash
	transferFnSignature := []byte("transfer(address,uint256)") // ERC-20 Smartcontrct Function at transfer
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	// Get MethodID
	methodID := hash.Sum(nil)[:4]                           //0xa9059cbb
	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)    // The left byte of the recipient hex address is 32 bytes filled with 0
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32) // The left byte of the amount is 32 bytes filled with 0
	var data []byte
	data = append(data, methodID...)      // transferFnSignature
	data = append(data, paddedAddress...) // Hex recipient address
	data = append(data, paddedAmount...)  // Hex amount
	return data
}

// SendNew1 godoc
func SendNew1(privateKeyHex, toStr string) (string, error) {

	value := big.NewInt(110000000000000) // in wei (0.00011 eth)
	gasLimit := uint64(30000000)         // in units
	tip := big.NewInt(2000000000)        // maxPriorityFeePerGas = 2 Gwei
	feeCap := big.NewInt(6000000000)     // maxFeePerGas = 6 Gwei
	//20000000000
	//9424646343000 = 4487926830*2100 = 0.000009424646343//eth //GasPrice

	// Connect the eth client to the nodeURL.
	client, err := ethClient()
	if err != nil {
		return "", err
	}

	// Toaddress
	ToAddress := common.HexToAddress(toStr)
	// Send Transfer signer From signer Private Key String
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	// Get ChainId from Ethereum Network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	// Get Nonce from From address
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return "", err
	}

	txData := []byte{} // Transfer ERC-20, so Eth Send not used

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,  //11155111
		Nonce:     nonce,    //67
		GasTipCap: tip,      //a.k.a. maxPriorityFeePerGas //2
		GasFeeCap: feeCap,   // a.k.a. maxFeePerGas //6
		Gas:       gasLimit, //300000
		To:        &ToAddress,
		Value:     value,
		Data:      txData,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// SendNew2 godoc
func SendNew2(ContractStr, privateKeyHex, toStr string) (string, error) {

	value := big.NewInt(0)            // in wei (0 eth)
	gasLimit := uint64(100000)        // in units
	tip := big.NewInt(1500000000)     // maxPriorityFeePerGas = 2 Gwei
	feeCap := big.NewInt(20000000000) // maxFeePerGas = 30 Gwei
	//20000000000
	//9424646343000 = 4487926830*2100 = 0.000009424646343//eth //GasPrice
	token := big.NewInt(0)
	token.SetUint64(1e18)

	// Connect the eth client to the nodeURL.
	client, err := ethClient()
	if err != nil {
		return "", err
	}

	// Toaddress
	ToAddress := common.HexToAddress(toStr)
	ContractAddress := common.HexToAddress(ContractStr)
	// Send Transfer signer From signer Private Key String
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	// Get ChainId from Ethereum Network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	// Get Nonce from From address
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return "", err
	}

	txData := SetTokenData(ToAddress, token) // Transfer ERC-20, so Eth Send not used

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,  //11155111
		Nonce:     nonce,    //67
		GasTipCap: tip,      //a.k.a. maxPriorityFeePerGas //2
		GasFeeCap: feeCap,   // a.k.a. maxFeePerGas //6
		Gas:       gasLimit, //300000
		To:        &ContractAddress,
		Value:     value,
		Data:      txData,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// SendNew3 godoc
func SendNew3(privateKeyHex, toStr string) (string, error) {

	value := big.NewInt(1040030000000000) // in wei (0.00104903 eth)
	gasLimit := uint64(300000)            // in units
	// tip := big.NewInt(2000000000)         // maxPriorityFeePerGas = 2 Gwei
	// feeCap := big.NewInt(6000000000)      // maxFeePerGas = 6 Gwei
	//20000000000
	//9424646343000 = 4487926830*2100 = 0.000009424646343//eth //GasPrice

	// Connect the eth client to the nodeURL.
	client, err := ethClient()
	if err != nil {
		return "", err
	}
	//SuggestGasPrice //22513306997  Wei
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	fmt.Printf("gasPrice1: %v\n", gasPrice)
	gasPrice.Add(gasPrice, big.NewInt(500000000)) // 0.5Gwei Added
	fmt.Printf("gasPrice2: %v\n", gasPrice)

	// Toaddress
	ToAddress := common.HexToAddress(toStr)
	// Send Transfer signer From signer Private Key String
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	// Get ChainId from Ethereum Network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	// Get Nonce from From address
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return "", err
	}

	txData := []byte{} // Transfer ERC-20, so Eth Send not used

	tx := types.NewTransaction(nonce, ToAddress, value, gasLimit, gasPrice, txData)
	//tx := types.NewContractCreation(nonce, value, gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// SendNew4 godoc
func SendNew4(ContractStr, privateKeyHex, toStr string) (string, error) {

	value := big.NewInt(0)     // in wei (0 eth)
	gasLimit := uint64(100000) // in units
	// tip := big.NewInt(1500000000)     // maxPriorityFeePerGas = 2 Gwei
	// feeCap := big.NewInt(20000000000) // maxFeePerGas = 30 Gwei
	//20000000000
	//9424646343000 = 4487926830*2100 = 0.000009424646343//eth //GasPrice
	token := big.NewInt(0)
	token.SetUint64(1e18)

	// Connect the eth client to the nodeURL.
	client, err := ethClient()
	if err != nil {
		return "", err
	}

	//SuggestGasPrice //22513306997  Wei
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	fmt.Printf("gasPrice1: %v\n", gasPrice)
	gasPrice.Add(gasPrice, big.NewInt(500000000)) // 0.5Gwei Added
	fmt.Printf("gasPrice2: %v\n", gasPrice)
	//gasPrice := big.NewInt(51000000000) //0.000000051//0.0040PHP
	//Transfer Fee 34,326*51000000000 //1,750,626,000,000,000// 0.001750626 eth //156.88PHP
	// Toaddress
	ToAddress := common.HexToAddress(toStr)
	ContractAddress := common.HexToAddress(ContractStr)
	// Send Transfer signer From signer Private Key String
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	// Get ChainId from Ethereum Network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	// Get Nonce from From address
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return "", err
	}

	txData := SetTokenData(ToAddress, token) // Transfer ERC-20, so Eth Send not used

	tx := types.NewTransaction(nonce, ContractAddress, value, gasLimit, gasPrice, txData)
	//tx := types.NewContractCreation(nonce, value, gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	//signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
