package main

import (
    "context"
    "fmt"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/mats9693/study/points/eth/utils"
    "log"
    "strings"
    "sync"
)

const (
    ethNodeAddr  = "ws://192.168.2.57:8546"
    contractAddr = "0x8a8CcB904ECf84B4ea43AEfCA1e6847640fB7f4c"

    contractABIStr = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"admin","type":"address"},{"indexed":false,"internalType":"uint8","name":"invokeTimes","type":"uint8"}],"name":"TestEventName","type":"event"},{"inputs":[],"name":"testFuncEmitEvent","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

var (
    ethConn *ethclient.Client
    sub     ethereum.Subscription

    contractABI  abi.ABI
    eventNameMap sync.Map // event hash - event name

    query = ethereum.FilterQuery{
        Addresses: []common.Address{common.HexToAddress(contractAddr)},
    }

    logs = make(chan types.Log, 10000)

    err error
)

func main() {
    registerEvents()

    ethConn, err = ethclient.Dial(ethNodeAddr)
    utils.CheckError(err, "dial eth client failed")

    sub, err = ethConn.SubscribeFilterLogs(context.Background(), query, logs)
    utils.CheckError(err, "subscribe logs failed")

    for {
        select {
        case err = <-sub.Err():
            log.Fatalln(err) // exit
        case vLog := <-logs:
            // core data
            fmt.Printf("> --- Receive a new event: %s ---\n", matchEventName(vLog.Topics[0].Hex()))
            fmt.Println("> Log Data:(parsed)    ", parseContractEvents(vLog.Topics[0].Hex(), vLog.Data))
            fmt.Println("> Block Number:        ", vLog.BlockNumber)
            fmt.Println("> Block Hash:          ", vLog.BlockHash.Hex())
            fmt.Println("> Transaction Hash:    ", vLog.TxHash.Hex())
            fmt.Println()

            // some other data
            fmt.Println("> Contract Address:    ", vLog.Address.Hex())
            fmt.Println("> Contract Topics:     ", printEventTopics(vLog.Topics))
            fmt.Println("> Log Data:(origin)    ", hexutil.Encode(vLog.Data))
            fmt.Println("> Transaction Index:   ", vLog.TxIndex)
            fmt.Println("> Log Index in Block:  ", vLog.Index)
            fmt.Println("> Is Removed:          ", vLog.Removed)
            fmt.Println()
        }
    }
}

func registerEvents() {
    eventDeclaration := []byte("TestEventName(address,uint8)")
    hash := crypto.Keccak256Hash(eventDeclaration)
    //fmt.Println("> Contract Event Signature(calc): ", hash.Hex() == "0x37bf82b399445377adc74da9876029ab2e1a0de7fedb054ecbf811afb4f6abe5", hash.Hex())

    eventNameMap.Store(hash.Hex(), "TestEventName")

    contractABI, err = abi.JSON(strings.NewReader(contractABIStr))
    utils.CheckError(err, "parse contract abi failed")
}

func matchEventName(hash string) string {
    eventNameI, ok := eventNameMap.Load(hash)
    if !ok {
        return ""
    }

    eventName, ok := eventNameI.(string)
    if !ok {
        return ""
    }

    return eventName
}

func printEventTopics(topics []common.Hash) []string {
    res := make([]string, 0, len(topics))
    for _, t := range topics {
        res = append(res, t.Hex())
    }

    return res
}

func parseContractEvents(hash string, data []byte) (res string) {
    switch hash {
    case "0x37bf82b399445377adc74da9876029ab2e1a0de7fedb054ecbf811afb4f6abe5": // TestEventName
        var payload []interface{}
        payload, err = contractABI.Unpack("TestEventName", data)
        utils.CheckError(err, "unpack params on event: TestEventName failed")
        utils.CheckBool(len(payload) == 2, "unexpected params amount")

        event := &EventParams_TestEventName{}

        paramOne, ok := payload[0].(common.Address)
        utils.CheckBool(ok, "type assert failed on event: TestEventName, param index: 0")
        event.Admin = paramOne

        paramTwo, ok := payload[1].(uint8)
        utils.CheckBool(ok, "type assert failed on event: TestEventName, param index: 1")
        event.InvokeTimes = paramTwo

        res = event.String()
    default:
        res = "unknown event hash: " + hash
        log.Println(res)
    }

    return
}
