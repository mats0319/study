package test

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/mats0319/secure_transfer/internal"
)

func TestEncDecOnce(t *testing.T) {
	err := internal.GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}

	err = internal.InitMessageFile()
	if err != nil {
		t.Error(err)
	}

	err = internal.Encrypt()
	if err != nil {
		t.Error(err)
	}

	t.Log("> Encrypt Success.")

	err = internal.Decrypt()
	if err != nil {
		t.Error(err)
	}

	t.Log("> Decrypt Success.")

	isEqual, err := compareFile(t, "message.txt", "message_decrypted.txt")
	if err != nil || !isEqual {
		t.Error(err)
	}

	t.Log("> File Hash Matched, Test Passed.")
}

func compareFile(t *testing.T, fileNameOne string, fileNameTwo string) (isEqual bool, err error) {
	fileOne, err := os.Open(fileNameOne)
	if err != nil {
		return
	}
	defer fileOne.Close()

	fileTwo, err := os.Open(fileNameTwo)
	if err != nil {
		return
	}
	defer fileTwo.Close()

	fileOneInfo, _ := fileOne.Stat()
	fileOneSize := fileOneInfo.Size()

	fileTwoInfo, _ := fileTwo.Stat()
	fileTwoSize := fileTwoInfo.Size()

	// compare file size
	if fileOneSize != fileTwoSize {
		return
	}

	hashOne, err := calcFileHash(fileOne)
	if err != nil {
		return
	}
	t.Log("- Origin    File Hash: ", hashOne)

	hashTwo, err := calcFileHash(fileTwo)
	if err != nil {
		return
	}
	t.Log("- Decrypted File Hash: ", hashTwo)

	isEqual = hashOne == hashTwo

	return
}

func calcFileHash(file *os.File) (hash string, err error) {
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return
	}

	hashBytes := hasher.Sum(nil)
	hash = hex.EncodeToString(hashBytes)

	return
}
