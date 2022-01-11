package utils

import (
    "context"
    "crypto/ecdsa"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "log"
    "math/big"
)

type TxParams struct {
    From    common.Address
    PrivKey *ecdsa.PrivateKey

    Signer types.Signer

    Value    *big.Int
    Pending  bool
    GasPrice *big.Int
    GasLimit uint64
}

func (t *TxParams) BuildCallOpts() *bind.CallOpts {
    return &bind.CallOpts{
        Pending:     t.Pending,
        From:        t.From,
        BlockNumber: nil,
        Context:     context.Background(),
    }
}

func (t *TxParams) BuildTxOpts() *bind.TransactOpts {
    return &bind.TransactOpts{
        From:  t.From,
        Nonce: nil,
        Signer: func(signer types.Signer, address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
            return types.SignTx(transaction, t.Signer, t.PrivKey)
        },
        Value:    t.Value,
        GasPrice: t.GasPrice,
        GasLimit: t.GasLimit,
        Context:  context.Background(),
    }
}

func NewTxParam(from, privKey string, chainID int64, gasPrice *big.Int, gasLimit uint64) *TxParams {
    pkey, err := crypto.HexToECDSA(privKey)
    if err != nil {
        log.Fatalln(err, "parse secp256k1 private key failed")
    }

    return &TxParams{
        From:    common.HexToAddress(from),
        PrivKey: pkey,

        Signer: types.NewEIP155Signer(big.NewInt(chainID)),

        Value:   big.NewInt(0),
        Pending: false,

        GasPrice: gasPrice,
        GasLimit: gasLimit,
    }
}
