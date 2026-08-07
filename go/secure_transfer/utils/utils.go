package utils

import (
	"crypto/ecdh"
	"math/rand/v2"
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

		if strings.HasPrefix(fileInfo.Name(), fileBaseName) { // match file name without extension
			fileName = fileInfo.Name()
			fileSize = fileInfo.Size()
			break
		}
	}

	return
}

const charactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示全部字符库中的字符
const useMask = 1<<useBits - 1

func GenerateRandomBytes[T string | []byte](length int) T {
	b := make([]byte, length)

	randomNum, remainBits := rand.Int64(), 64
	for i := 0; i < len(b); {
		if remainBits < useBits {
			randomNum, remainBits = rand.Int64(), 64
		}

		index := int(randomNum & useMask) // 0b0011 1111
		if index < len(charactersLibrary) {
			randomNum >>= useBits
			remainBits -= useBits

			b[i] = charactersLibrary[index]
			i++
		} else {
			randomNum >>= 1
			remainBits -= 1
		}
	}

	return T(b)
}
