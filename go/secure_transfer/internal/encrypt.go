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
)

func Encrypt() (isSuccess bool) {
	pubKey := deserializePublicKey()
	if pubKey == nil {
		return
	}

	fileName, fileBytes := GetFirstFile(plainTextFileName)
	if len(fileBytes) < 1 {
		return
	}

	encryptedBytes := encrypt(pubKey, fileBytes)
	if encryptedBytes == nil {
		return
	}

	extension := strings.ToUpper(GetExtension(fileName, plainTextFileName))
	err := os.WriteFile(fmt.Sprintf("./%s%s", cipherFileName, extension), encryptedBytes, 0777)
	if err != nil {
		Error("Write cipher file", err)
		return
	}

	Success("Encrypt")

	return true
}

func deserializePublicKey() *ecdh.PublicKey {
	pubKeyBytes, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		Error("Read public key", err)
		return nil
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		Error("Decode public key", nil)
		return nil
	}

	pubKeyI, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		Error("Parse public key", err)
		return nil
	}

	pubKey, ok := pubKeyI.(*ecdh.PublicKey)
	if !ok {
		ecdsaPubKey, ok := pubKeyI.(*ecdsa.PublicKey) // x509 parser usually return this type
		if !ok {
			Error("Public key type assert", nil)
			return nil
		}

		pubKey, err = ecdsaPubKey.ECDH()
		if err != nil {
			Error("ECDSA public key to ecdh", err)
			return nil
		}
	}

	return pubKey
}

func encrypt(pubKey *ecdh.PublicKey, content []byte) []byte {
	pubKeyBytes, aesKey, err := deriveKeyInEncrypt(pubKey)
	if err != nil {
		return nil
	}

	ciphertext := aesEncrypt(aesKey, content)

	result := make([]byte, 1+len(pubKeyBytes)+len(ciphertext))
	copy(result[:1], []byte{byte(len(pubKeyBytes))})
	copy(result[1:1+len(pubKeyBytes)], pubKeyBytes)
	copy(result[1+len(pubKeyBytes):], ciphertext)

	return result
}

func deriveKeyInEncrypt(pubKey *ecdh.PublicKey) (pubKeyBytes []byte, aesKey []byte, err error) {
	// ecdh
	tempPrivKey, err := Curve().GenerateKey(nil)
	if err != nil {
		Error("Generate temp key", err)
		return
	}

	pubKeyBytes = tempPrivKey.PublicKey().Bytes()

	sharedKey, err := tempPrivKey.ECDH(pubKey)
	if err != nil {
		Error("ECDH", err)
		return
	}

	// kdf
	aesKey, err = hkdf.Key(sha256.New, sharedKey, []byte(salt), info, driveKeyLength)
	if err != nil {
		Error("HKDF", err)
		return
	}

	return
}

// aes-gcm encrypt
func aesEncrypt(aesKey []byte, content []byte) []byte {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		Error("Build cipher block", err)
		return nil
	}

	aesGCM, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		Error("Build GCM", err)
		return nil
	}

	return aesGCM.Seal(nil, nil, content, nil)
}
