package main

import (
    "context"
    "fmt"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/mats9693/study/points/eth/utils"
    "strings"
)

const (
    ethNodeAddr = "http://192.168.2.57:8545"

    address = "0x552504ef1e81786ccf974c8ef27bfdb2c12e7b86"

    contractAddressStr = "0x713EDF08219c4922878ce91A62622d121CDE7D36"
    contractABIStr     = `[{"inputs":[{"internalType":"address","name":"erc721ContractAddress","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"string","name":"extendedParam","type":"string"}],"name":"AuctionLock","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"string","name":"extendedParam","type":"string"}],"name":"AuctionUnlock","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"string","name":"batchCode","type":"string"},{"indexed":false,"internalType":"uint16","name":"tokenAmount","type":"uint16"},{"indexed":false,"internalType":"uint16[]","name":"tokenTypes","type":"uint16[]"},{"indexed":false,"internalType":"uint256[]","name":"tokenIds","type":"uint256[]"}],"name":"Mint","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"address","name":"buyer","type":"address"},{"indexed":false,"internalType":"string","name":"extendedParam","type":"string"}],"name":"OfficialAuctionSucceed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"string","name":"batchCode","type":"string"},{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"uint16","name":"tokenType","type":"uint16"}],"name":"OpenBlindBox","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"string","name":"extendedParam","type":"string"}],"name":"RechargeBCash","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokenId","type":"uint256"},{"indexed":false,"internalType":"address","name":"buyer","type":"address"},{"indexed":false,"internalType":"string","name":"extendedParam","type":"string"}],"name":"UserAuctionSucceed","type":"event"},{"inputs":[{"internalType":"string","name":"batchCode","type":"string"},{"internalType":"uint16[]","name":"types","type":"uint16[]"},{"internalType":"uint256[]","name":"tokenIds","type":"uint256[]"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"address","name":"buyer","type":"address"},{"internalType":"string","name":"extendedParam","type":"string"}],"name":"officialAuctionSucceed","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"batchCode","type":"string"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"uint16","name":"tokenType","type":"uint16"}],"name":"openBlindBox","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"string","name":"extendedParam","type":"string"}],"name":"rechargeBCash","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"string","name":"extendedParam","type":"string"}],"name":"auctionLock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"string","name":"extendedParam","type":"string"}],"name":"auctionUnlock","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"address","name":"buyer","type":"address"},{"internalType":"string","name":"extendedParam","type":"string"}],"name":"userAuctionSucceed","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newAddr","type":"address"}],"name":"setAdmin","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"newOperators","type":"address[]"}],"name":"setOperators","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"newMiners","type":"address[]"}],"name":"setMiners","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getOperators","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getMiners","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"batchCode","type":"string"}],"name":"getBatchTokens","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"batchCode","type":"string"}],"name":"getBatchTypes","outputs":[{"internalType":"uint16[]","name":"","type":"uint16[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"batchCode","type":"string"}],"name":"getBatchOpenStatus","outputs":[{"internalType":"uint16[]","name":"","type":"uint16[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"getTokenStatus","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"getTokenType","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"stateMutability":"view","type":"function"}]`

    contractMethodName = "getMiners"
)

var (
    ethConn         *ethclient.Client
    contractAddress common.Address
    contractABI     abi.ABI
    txData          []byte

    err error
)

func main() {
    ethConn, err = ethclient.Dial(ethNodeAddr)
    utils.CheckError(err, "dial eth client failed")

    contractAddress = common.HexToAddress(contractAddressStr)

    contractABI, err = abi.JSON(strings.NewReader(contractABIStr))
    utils.CheckError(err, "parse contract abi failed")

    txData, err = contractABI.Pack(contractMethodName)
    utils.CheckError(err, "pack params failed")

    getMiners()
}

func getMiners() {
    var res []byte
    res, err = ethConn.CallContract(context.Background(), ethereum.CallMsg{
        From: common.HexToAddress(address),
        To:   &contractAddress,
        Data: txData,
    }, nil)
    utils.CheckError(err, "contract call failed")

    fmt.Println("> res: ", hexutil.Encode(res))
}
