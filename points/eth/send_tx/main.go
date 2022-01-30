package main

import (
    "context"
    "crypto/ecdsa"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/mats9693/study/points/eth/utils"
    "log"
    "math/big"
    "strings"
)

const (
    ethNodeAddr            = "http://192.168.2.57:8545"
    defaultGasLimit uint64 = 5000000
    defaultGasPrice        = 0
    chainID                = 9876

    address    = "0xe95534f7d843f71873c768049b5fc5cbbb850ec8"
    privKeyStr = "c12e899402274bda6231a35d509a7b15bf69dba64f125b7de4655d3b5eb9139a"

    // in truffle, the long code
    contractByteCodeForDeploy = "0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101fd806100606000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80632973d0be14610030575b600080fd5b61003861003a565b005b6001600060148282829054906101000a900460ff166100599190610122565b92506101000a81548160ff021916908360ff1602179055507f37bf82b399445377adc74da9876029ab2e1a0de7fedb054ecbf811afb4f6abe560008054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600060149054906101000a900460ff166040516100d19291906100f9565b60405180910390a1565b6100e481610159565b82525050565b6100f38161018b565b82525050565b600060408201905061010e60008301856100db565b61011b60208301846100ea565b9392505050565b600061012d8261018b565b91506101388361018b565b92508260ff0382111561014e5761014d610198565b5b828201905092915050565b60006101648261016b565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600060ff82169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea26469706673582212206b3d2a4c2cb57bd7d09ef4018209eb7ea9ef9fc5e33b40c25843c8cb6085159864736f6c63430008060033"

    contractAddressStr = "0x8a8CcB904ECf84B4ea43AEfCA1e6847640fB7f4c"
    contractABIStr     = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"admin","type":"address"},{"indexed":false,"internalType":"uint8","name":"invokeTimes","type":"uint8"}],"name":"TestEventName","type":"event"},{"inputs":[],"name":"testFuncEmitEvent","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
    contractMethodName = "testFuncEmitEvent"
)

var (
    ethConn *ethclient.Client

    nonce   uint64
    privKey *ecdsa.PrivateKey

    tx       *types.Transaction
    signedTx *types.Transaction

    contractABI     abi.ABI
    contractAddress = common.HexToAddress(contractAddressStr)
    txData          []byte

    err error
)

func main() {
    // connect eth client and prepare nonce, private key
    {
        ethConn, err = ethclient.Dial(ethNodeAddr)
        utils.CheckError(err, "dial eth client failed")

        nonce, err = ethConn.NonceAt(context.Background(), common.HexToAddress(address), nil)
        utils.CheckError(err, "get nonce failed")

        privKey, err = crypto.HexToECDSA(privKeyStr)
        utils.CheckError(err, "parse secp256k1 private key failed")
    }

    // test connect, get last block number
    //getLastBlockHeader()

    // deploy contract
    //deployContract()
    //getContractAddress()

    // invoke contract function
    invokeContractFunction()
}

func getLastBlockHeader() {
    header, err := ethConn.HeaderByNumber(context.Background(), nil)
    utils.CheckError(err, "get last block header failed")

    log.Println("> Last block: ", header.Number.Uint64())
}

func deployContract() {
    // build tx and sign
    tx = types.NewContractCreation(nonce, nil, defaultGasLimit, nil, common.FromHex(contractByteCodeForDeploy))
    signedTx, err = types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privKey)
    utils.CheckError(err, "sign tx failed")

    log.Println("> ------- calc -------")
    log.Println("> tx hash :", signedTx.Hash().Hex())
    log.Println("> tx cost :", signedTx.Cost().String())
    log.Println("> tx nonce:", signedTx.Nonce())

    // send tx
    err = ethConn.SendTransaction(context.Background(), signedTx)
    utils.CheckError(err, "send tx failed")
    log.Println("> success.")
}

func getContractAddress() {
    var receipt *types.Receipt
    receipt, err = ethConn.TransactionReceipt(context.Background(), common.HexToHash("0x94d27171a46303acaa703c71455729226c6e46b2e4b8e532a9204e02be86d2c3"))
    utils.CheckError(err, "get receipt failed")

    contractAddrCalc := crypto.CreateAddress(common.HexToAddress("0xe95534f7d843f71873c768049b5fc5cbbb850ec8"), 2)

    log.Println("> Contract address      : ", receipt.ContractAddress.Hex())
    log.Println("> Contract address(calc): ", contractAddrCalc.Hex())
}

func invokeContractFunction() {
    // build tx data according to contract ABI
    contractABI, err = abi.JSON(strings.NewReader(contractABIStr))
    utils.CheckError(err, "parse contract abi failed")

    txData, err = contractABI.Pack(contractMethodName)
    utils.CheckError(err, "pack params failed")

    // pre-execute tx
    _, err = ethConn.EstimateGas(context.Background(), ethereum.CallMsg{
        From:      common.HexToAddress(address),
        To:        &contractAddress,
        GasFeeCap: big.NewInt(defaultGasPrice),
        GasTipCap: big.NewInt(defaultGasPrice),
        Data:      txData,
    })
    utils.CheckError(err, "tx exec failed")

    // build tx and sign
    tx = types.NewTx(&types.LegacyTx{
        Nonce:    nonce,
        GasPrice: big.NewInt(defaultGasPrice),
        Gas:      defaultGasLimit,
        To:       &contractAddress,
        Value:    big.NewInt(0),
        Data:     txData,
    })

    signedTx, err = types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privKey)
    utils.CheckError(err, "sign tx failed")

    // send tx
    err = ethConn.SendTransaction(context.Background(), signedTx)
    utils.CheckError(err, "send tx failed")
    log.Println("> success. ", signedTx.Hash().Hex())
}
