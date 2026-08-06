package utils

import (
	"crypto/ecdh"
	"os"
	"strings"

	mlog "github.com/mats0319/secure_transfer/utils/log"
)

// Curve use same curve unified
func Curve() ecdh.Curve {
	return ecdh.X25519()
}

// FirstFile fileBaseName: file name without extension
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

		if strings.HasPrefix(fileInfo.Name(), fileBaseName) { // match file name without extension
			fileName = fileInfo.Name()
			fileSize = fileInfo.Size()
			break
		}
	}

	return
}

func GetFirstFile(fileName string) (filePath string, fileBytes []byte, e *Error) {
	entry, err := os.ReadDir("./")
	if err != nil {
		e = ErrReadDir().WithCause(err)
		mlog.Error(e.String())
		return
	}

	var filePathBuilder strings.Builder
	filePathBuilder.WriteString("./")

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

		if strings.HasPrefix(fileInfo.Name(), fileName) { // match file name without extension
			filePathBuilder.WriteString(fileInfo.Name())
			break
		}
	}

	filePath = filePathBuilder.String()

	fileBytes, err = os.ReadFile(filePath)
	if err != nil {
		e = ErrReadMessageFile().WithCause(err).WithParam("name", filePath)
		mlog.Error(e.String())
		return
	}

	return
}

func GetExtension(filePath string, fileName string) string {
	return strings.TrimPrefix(filePath, "./"+fileName)
}
