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
	decFileName := fmt.Sprintf("%s.%s", utils.DecryptedFileName, strings.ToUpper(extension))

	// select decryptor
	switch fileHeader.EncMethod {
	case utils.EncryptMethod_Once:
		dec := lib.NewDecryptorOnce(privKey, fileHeader, encFileSize)
		e := dec.Decrypt(encFileName, decFileName)
		if e != nil {
			return e
		}
	case utils.EncryptMethod_Stream:
		// todo dec stream
	}

	return nil
}

func deserializePrivateKey() (privKey *ecdh.PrivateKey, e *utils.Error) {
	privKeyBytes, err := os.ReadFile(privateKeyFilePath)
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
			e = utils.ErrInvalidPrivateKey()
			mlog.Error(e.String())
			return
		}

		privKey, err = ecdsaPrivKey.ECDH()
		if err != nil {
			e = utils.ErrInvalidPrivateKey().WithCause(err)
			mlog.Error(e.String())
			return
		}
	}

	return
}

func decrypt(privKey *ecdh.PrivateKey, pubKeyFileBytes []byte) (decryptedBytes []byte, e *utils.Error) {
	aesKey, ciphertext, e := deriveKeyInDecrypt(privKey, pubKeyFileBytes)
	if e != nil {
		return
	}

	decryptedBytes, e = aesDecrypt(aesKey, ciphertext)
	if e != nil {
		return
	}

	return
}

func deriveKeyInDecrypt(privKey *ecdh.PrivateKey, pubKeyFileBytes []byte) (aesKey []byte, ciphertext []byte, e *utils.Error) {
	// ecdh
	if len(pubKeyFileBytes) <= 1+publicKeyLength || pubKeyFileBytes[0] != publicKeyLength {
		e = utils.ErrInvalidPublicKey().WithParam("length", len(pubKeyFileBytes))
		mlog.Error(e.String())
		return
	}

	pubKeyBytes := pubKeyFileBytes[1 : 1+publicKeyLength]
	ciphertext = pubKeyFileBytes[1+publicKeyLength:]

	pubKey, err := utils.Curve().NewPublicKey(pubKeyBytes)
	if err != nil {
		e = utils.ErrInvalidPublicKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	sharedKey, err := privKey.ECDH(pubKey)
	if err != nil {
		e = utils.ErrECDH().WithCause(err)
		mlog.Error(e.String())
		return
	}

	// kdf
	aesKey, err = hkdf.Key(sha256.New, sharedKey, []byte(salt), info, deriveKeyLength)
	if err != nil {
		e = utils.ErrKDF().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}

func aesDecrypt(aesKey []byte, ciphertext []byte) (decryptedBytes []byte, e *utils.Error) {
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

	decryptedBytes, err = aesGCM.Open(nil, nil, ciphertext, nil)
	if err != nil {
		e = utils.ErrAESDecrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
