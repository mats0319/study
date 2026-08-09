package test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

// fileSize: unit 1m
func generateRandomFile(fileName string, fileSize int) error {
	return exec.Command("dd",
		"if=/dev/urandom",
		"of="+fileName,
		fmt.Sprintf("bs=%dM", fileSize),
		"count=1",
	).Run()
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
