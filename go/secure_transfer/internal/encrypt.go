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

func Encrypt() error {
	pubKey, e := deserializePublicKey()
	if e != nil {
		return e
	}

	// get file name(s)
	originFileName, originFileSize, e := utils.FirstFile(utils.OriginFileName)
	if e != nil {
		return e
	}

	index := strings.LastIndex(originFileName, ".")
	extension := originFileName[index+1:]
	encFileName := fmt.Sprintf("%s.%s", utils.EncryptedFileName, strings.ToUpper(extension))

	// select encryptor
	var enc lib.Encryptor
	if originFileSize <= utils.OnceEncMaxSize {
		enc = lib.NewEncryptorOnce(pubKey)
	} else {
		enc = lib.NewEncryptorStream(pubKey)
	}

	e = enc.Encrypt(originFileName, encFileName)
	if e != nil {
		_ = os.Remove(encFileName) // 删除可能存在的半截文件
		return e
	}

	return nil
}

func deserializePublicKey() (pubKey *ecdh.PublicKey, e *utils.Error) {
	pubKeyBytes, err := os.ReadFile(utils.PublicKeyFileName)
	if err != nil {
		e = utils.ErrReadPublicKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		e = utils.ErrDecodePublicKey()
		mlog.Error(e.String())
		return
	}

	pubKeyI, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		e = utils.ErrParsePublicKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	pubKey, ok := pubKeyI.(*ecdh.PublicKey)
	if !ok {
		ecdsaPubKey, ok := pubKeyI.(*ecdsa.PublicKey) // x509 parser usually return this type
		if !ok {
			e = utils.ErrParsePublicKey()
			mlog.Error(e.String())
			return
		}

		pubKey, err = ecdsaPubKey.ECDH()
		if err != nil {
			e = utils.ErrParsePublicKey().WithCause(err)
			mlog.Error(e.String())
			return
		}
	}

	return
}
