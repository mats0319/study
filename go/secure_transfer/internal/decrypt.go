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
	"errors"
	"fmt"
	"os"
	"strings"
)

func Decrypt() (isSuccess bool) {
	privKey := deserializePrivateKey()
	if privKey == nil {
		return
	}

	fileName, fileBytes := GetFirstFile(cipherFileName)
	if len(fileBytes) < 1 {
		return
	}

	decryptedBytes := decrypt(privKey, fileBytes)
	if decryptedBytes == nil {
		return
	}

	extension := strings.ToLower(GetExtension(fileName, cipherFileName))
	err := os.WriteFile(fmt.Sprintf("./%s%s", plainTextDecryptedFileName, extension), decryptedBytes, 0777)
	if err != nil {
		Error("Write message", err)
		return
	}

	Success("Decrypt")

	return true
}

func deserializePrivateKey() *ecdh.PrivateKey {
	privKeyBytes, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		Error("Read private key", err)
		return nil
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		Error("Decode private key", nil)
		return nil
	}

	privKeyI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Error("Parse private key", err)
		return nil
	}

	privKey, ok := privKeyI.(*ecdh.PrivateKey)
	if !ok {
		ecdsaPrivKey, ok := privKeyI.(*ecdsa.PrivateKey)
		if !ok {
			Error("Private key type assert", nil)
			return nil
		}

		privKey, err = ecdsaPrivKey.ECDH()
		if err != nil {
			Error("ECDSA private key to ecdh", err)
			return nil
		}
	}

	return privKey
}

func decrypt(privKey *ecdh.PrivateKey, pubKeyFileBytes []byte) []byte {
	aesKey, ciphertext, err := deriveKeyInDecrypt(privKey, pubKeyFileBytes)
	if err != nil {
		return nil
	}

	message, err := aesDecrypt(aesKey, ciphertext)
	if err != nil {
		return nil
	}

	return message
}

func deriveKeyInDecrypt(privKey *ecdh.PrivateKey, pubKeyFileBytes []byte) (aesKey []byte, ciphertext []byte, err error) {
	if len(pubKeyFileBytes) <= 1+publicKeyLength || pubKeyFileBytes[0] != publicKeyLength {
		err = errors.New(fmt.Sprintf("wrong length, want: %d, get: %d", publicKeyLength, len(pubKeyFileBytes)))
		Error("Invalid ciphertext", err)
		return
	}

	pubKeyBytes := pubKeyFileBytes[1 : 1+publicKeyLength]
	ciphertext = pubKeyFileBytes[1+publicKeyLength:]

	pubKey, err := Curve().NewPublicKey(pubKeyBytes)
	if err != nil {
		Error("Load public key", err)
		return
	}

	sharedKey, err := privKey.ECDH(pubKey)
	if err != nil {
		Error("ECDH", err)
		return
	}

	aesKey, err = hkdf.Key(sha256.New, sharedKey, []byte(salt), info, driveKeyLength)
	if err != nil {
		Error("HKDF", err)
		return
	}

	return
}

func aesDecrypt(aesKey []byte, ciphertext []byte) (message []byte, err error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		Error("Build cipher block", err)
		return
	}

	aesGCM, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		Error("Build GCM", err)
		return
	}

	message, err = aesGCM.Open(nil, nil, ciphertext, nil)
	if err != nil {
		Error("Decrypt", err)
		return
	}

	return
}
