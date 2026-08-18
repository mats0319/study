package internal

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/mats0319/secure_transfer/internal/lib"
	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

func Decrypt() error {
	privKey, e := deserializePrivateKey()
	if e != nil {
		return e
	}

	// get file name(s)
	encFileName, encFileSize, e := utils.FirstFile(utils.EncryptedFileName)
	if e != nil {
		return e
	}

	fileHeader := &lib.FileHeader{}
	e = fileHeader.Deserialize(encFileName)
	if e != nil {
		return e
	}

	index := strings.LastIndex(encFileName, ".")
	extension := encFileName[index+1:]
	decFileName := fmt.Sprintf("%s.%s", utils.DecryptedFileName, strings.ToLower(extension))

	// select decryptor
	var dec lib.Decryptor
	switch fileHeader.EncMethod {
	case utils.EncryptMethod_Once:
		dec = lib.NewDecryptorOnce(privKey, fileHeader, encFileSize)
	case utils.EncryptMethod_Stream:
		dec = lib.NewDecryptorStream(privKey, fileHeader)
	}

	e = dec.Decrypt(encFileName, decFileName)
	if e != nil {
		_ = os.Remove(decFileName) // 删除可能存在的半截文件
		return e
	}

	return nil
}

func deserializePrivateKey() (privKey *ecdh.PrivateKey, e *utils.Error) {
	privKeyBytes, err := os.ReadFile(utils.PrivateKeyFileName)
	if err != nil {
		e = utils.ErrReadPrivateKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		e = utils.ErrDecodePrivateKey()
		mlog.Error(e.String())
		return
	}

	privKeyI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		e = utils.ErrParsePrivateKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	privKey, ok := privKeyI.(*ecdh.PrivateKey)
	if !ok {
		ecdsaPrivKey, ok := privKeyI.(*ecdsa.PrivateKey)
		if !ok {
			e = utils.ErrParsePrivateKey()
			mlog.Error(e.String())
			return
		}

		privKey, err = ecdsaPrivKey.ECDH()
		if err != nil {
			e = utils.ErrParsePrivateKey().WithCause(err)
			mlog.Error(e.String())
			return
		}
	}

	return
}
