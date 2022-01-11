package contract

import (
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/mats9693/study/points/eth/contract/stub"
    "github.com/mats9693/study/points/eth/utils"
    "github.com/pkg/errors"
    "math/big"
)

type contractWrapperImpl struct {
    data *contractinterface.Data
}

func (c *contractWrapperImpl) SetData(txParams *utils.TxParams, strs []string, nums []*big.Int) (*types.Transaction, error) {
    return c.data.SetData(txParams.BuildTxOpts(), strs, nums)
}

func (c *contractWrapperImpl) SetSingleData(txParams *utils.TxParams, str string, num *big.Int) (*types.Transaction, error) {
    return c.data.SetSingleData(txParams.BuildTxOpts(), str, num)
}

func NewContractWrapper(contractAddress common.Address, ethNodeAddr string) (ContractWrapper, error) {
    cwi := &contractWrapperImpl{}

    conn, err := ethclient.Dial(ethNodeAddr)
    if err != nil {
        return nil, errors.Wrap(err, fmt.Sprintf("connect to node: %s failed", ethNodeAddr))
    }

    cwi.data, err = contractinterface.NewData(contractAddress, conn)
    if err != nil {
        return nil, errors.Wrap(err, "init tokenx contract interface failed")
    }

    return cwi, nil
}
