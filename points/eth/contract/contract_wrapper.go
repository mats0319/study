package contract

import (
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/mats9693/study/points/eth/utils"
    "math/big"
)

// ContractWrapper 合约api抽象，将合约接口包装一层，
// 将 txParam -> callOpts/txOpts 的过程，以及eth连接对外隐藏
type ContractWrapper interface {
    SetData(txParams *utils.TxParams, strs []string, nums []*big.Int) (*types.Transaction, error)
    SetSingleData(txParmas *utils.TxParams, str string, num *big.Int) (*types.Transaction, error)
}
