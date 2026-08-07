package internal

import (
	"crypto/ecdh"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

func GenerateKeyPair() error {
	privKey, err := utils.Curve().GenerateKey(nil)
	if err != nil {
		e := utils.ErrGeneratePrivateKey().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	e := serializePrivateKey(privKey)
	if e != nil {
		return e
	}

	e = serializePublicKey(privKey.PublicKey())
	if e != nil {
		return e
	}

	return nil
}

func serializePrivateKey(privKey *ecdh.PrivateKey) *utils.Error {
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		e := utils.ErrMarshalPrivateKey().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	block := &pem.Block{Type: "Private Key", Bytes: privKeyBytes}
	blockBytes := pem.EncodeToMemory(block)

	err = os.WriteFile(utils.PrivateKeyFileName, blockBytes, 0600)
	if err != nil {
		e := utils.ErrSavePrivateKey().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func serializePublicKey(pubKey *ecdh.PublicKey) *utils.Error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		e := utils.ErrMarshalPublicKey().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	block := &pem.Block{Type: "Public Key", Bytes: pubKeyBytes}
	blockBytes := pem.EncodeToMemory(block)

	err = os.WriteFile(utils.PublicKeyFileName, blockBytes, 0644)
	if err != nil {
		e := utils.ErrSavePublicKey().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}
