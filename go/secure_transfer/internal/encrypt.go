package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/hkdf"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

func Encrypt() error {
	pubKey, e := deserializePublicKey()
	if e != nil {
		return e
	}

	filePath, fileBytes, e := utils.GetFirstFile(plainTextFileName)
	if e != nil {
		return e
	}

	encryptedBytes, e := encrypt(pubKey, fileBytes)
	if e != nil {
		return e
	}

	extension := strings.ToUpper(utils.GetExtension(filePath, plainTextFileName))
	fileName := fmt.Sprintf("./%s%s", cipherFileName, extension)
	err := os.WriteFile(fileName, encryptedBytes, 0777)
	if err != nil {
		e = utils.ErrSaveCiphertext().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func deserializePublicKey() (pubKey *ecdh.PublicKey, e *utils.Error) {
	pubKeyBytes, err := os.ReadFile(publicKeyFilePath)
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
			e = utils.ErrInvalidPublicKey()
			mlog.Error(e.String())
			return
		}

		pubKey, err = ecdsaPubKey.ECDH()
		if err != nil {
			e = utils.ErrInvalidPublicKey().WithCause(err)
			mlog.Error(e.String())
			return
		}
	}

	return
}

func encrypt(pubKey *ecdh.PublicKey, content []byte) (encryptedBytes []byte, e *utils.Error) {
	pubKeyBytes, aesKey, e := deriveKeyInEncrypt(pubKey)
	if e != nil {
		return
	}

	ciphertext, e := aesEncrypt(aesKey, content)
	if e != nil {
		return
	}

	encryptedBytes = make([]byte, 1+len(pubKeyBytes)+len(ciphertext))
	copy(encryptedBytes[:1], []byte{byte(len(pubKeyBytes))})
	copy(encryptedBytes[1:1+len(pubKeyBytes)], pubKeyBytes)
	copy(encryptedBytes[1+len(pubKeyBytes):], ciphertext)

	return
}

func deriveKeyInEncrypt(pubKey *ecdh.PublicKey) (pubKeyBytes []byte, aesKey []byte, e *utils.Error) {
	// ecdh
	tempPrivKey, err := utils.Curve().GenerateKey(nil)
	if err != nil {
		e = utils.ErrGeneratePrivateKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	pubKeyBytes = tempPrivKey.PublicKey().Bytes()

	sharedKey, err := tempPrivKey.ECDH(pubKey)
	if err != nil {
		e = utils.ErrECDH().WithCause(err)
		mlog.Error(e.String())
		return
	}

	// kdf
	aesKey, err = hkdf.Key(sha256.New, sharedKey, []byte(salt), info, driveKeyLength)
	if err != nil {
		e = utils.ErrKDF().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}

// aes-gcm encrypt
func aesEncrypt(aesKey []byte, content []byte) (encryptedBytes []byte, e *utils.Error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		e = utils.ErrAESNewCipher().WithCause(err)
		mlog.Error(e.String())
		return
	}

	aesGCM, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		e = utils.ErrAESNewGCM().WithCause(err)
		mlog.Error(e.String())
		return
	}

	encryptedBytes = aesGCM.Seal(nil, nil, content, nil)

	return
}
