package utils

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"os"
	"strings"

	mlog "github.com/mats0319/secure_transfer/utils/log"
)

// Curve use same curve unified
func Curve() ecdh.Curve {
	return ecdh.X25519()
}

// FirstFile get first file which name matched 'file base name'
// - fileBaseName: file name without extension
func FirstFile(fileBaseName string) (fileName string, fileSize int64, e *Error) {
	entry, err := os.ReadDir("./")
	if err != nil {
		e = ErrReadDir().WithCause(err)
		mlog.Error(e.String())
		return
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		fileInfo, err := entry[i].Info()
		if err != nil {
			e = ErrGetFileInfo().WithCause(err).WithParam("name", entry[i].Name())
			mlog.Error(e.String())
			continue
		}

		// match file name, require: 'message.xxx'
		if strings.HasPrefix(fileInfo.Name(), fileBaseName+".") &&
			(strings.Index(fileInfo.Name(), ".") == strings.LastIndex(fileInfo.Name(), ".")) {
			fileName = fileInfo.Name()
			fileSize = fileInfo.Size()
			break
		}
	}

	if fileName == "" {
		e = ErrNoMatchedFile()
	}

	return
}

func GenerateRandomBytes(length int) []byte {
	bytesBuilder := bytes.NewBuffer(nil)

	for l := 0; l < length; {
		n, _ := bytesBuilder.WriteString(rand.Text()) // err always nil
		l += n
	}

	return bytesBuilder.Bytes()[:length]
}

func CalcSHA256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)

	return hasher.Sum(nil)
}

func CompareBytes(a []byte, b []byte) (isEqual bool) {
	if len(a) != len(b) {
		return
	}

	isEqual = true
	for i := range a {
		if a[i] != b[i] {
			isEqual = false
			break
		}
	}

	return
}
