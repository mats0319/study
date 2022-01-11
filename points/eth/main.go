package main

import (
    "context"
    "crypto/ecdsa"
    "fmt"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/mats9693/study/points/eth/utils"
    "math/big"
    "strings"
    "time"
)

// eth connect
var (
    ethConn *ethclient.Client
    nonce   uint64

    tx        *types.Transaction
    signedTx  *types.Transaction
    daPrivKey *ecdsa.PrivateKey

    err error
)

// contract
var (
    contractAddress = "0x9186e87b17552b65bc0c1ebbd5d0811c2128c620"
    contractABI     abi.ABI
    txData          []byte
    strs1           []string
    strs2           []string
    nums            []*big.Int
)

// create contract
func main() {
    // make eth connection struct
    ethConn, err = ethclient.Dial(ethNodeAddr)
    utils.CheckError(err, "make eth connection structure failed, error:")

    // test: test connect (get balance)
    {
        balance, err := ethConn.BalanceAt(context.Background(), deployAddress, nil)
        utils.CheckError(err, "get balance failed, error:")
        utils.ShowParam(fmt.Sprintf("> balance of %s: %s\n", deployAddress.Hex(), balance.String()), "")
    }

    // prepare for deploy contract
    {
        // get nonce
        nonce, err = ethConn.NonceAt(context.Background(), deployAddress, nil)
        utils.CheckError(err, "get nonce failed, error:")
        //ShowParam("> nonce:", nonce)

        // make priv key from str
        daPrivKey, err = crypto.HexToECDSA(daPrivKeyStr)
        utils.CheckError(err, "parse secp256k1 private key failed, error:")
    }

    //deployContract()

    // prepare for invoke contract functions
    {
        makeContractInvokeData()

        contractABI, err = abi.JSON(strings.NewReader(contractABIStr))
        utils.CheckError(err, "read contract abi failed, error:")
    }

    {
        txData, err = contractABI.Pack("setData", strs1, nums)
        utils.CheckError(err, "pack params failed, error:")

        //txData, err = contractABI.Pack("setSingleData", strs2[0], nums[0])
        //CheckError(err, "pack params failed, error:")
    }

    //sendTx()
}

func deployContract() {
    // make tx and sign
    tx = types.NewContractCreation(nonce, big.NewInt(0), defaultGasLimit, defaultGasPrice, common.FromHex(contractByteCodeHex))
    signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), daPrivKey)
    utils.CheckError(err, "sign tx failed, error:")

    utils.ShowParam("> ------- calc -------", "")
    utils.ShowParam("> tx hash:", signedTx.Hash().Hex())
    utils.ShowParam("> tx cost:", signedTx.Cost().String())

    // send tx
    err = ethConn.SendTransaction(context.Background(), signedTx)
    utils.CheckError(err, "send tx failed, error:")
    utils.ShowParam("> success.", "")
}

func sendTx() {
    // make tx and sign
    tx = types.NewTransaction(nonce, common.HexToAddress(contractAddress), big.NewInt(0), defaultGasLimit, defaultGasPrice, txData)
    signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), daPrivKey)
    utils.CheckError(err, "sign tx failed, error:")

    // send tx
    err = ethConn.SendTransaction(context.Background(), signedTx)
    utils.CheckError(err, "send tx failed, error:")
    utils.ShowParam("> success.", signedTx.Hash().Hex())
}

func makeContractInvokeData() {
    timestamp := time.Now()
    length := 6
    for i := 0; i < length; i++ {
        if i < length/2 {
            strs1 = append(strs1, timestamp.Add(time.Second*time.Duration(i)).String())
            nums = append(nums, big.NewInt(timestamp.Unix()))
        } else {
            strs2 = append(strs2, timestamp.Add(time.Second*time.Duration(i)).String())
        }
    }
}
