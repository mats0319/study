package utils

import (
	"bytes"
	"crypto/rand"
)

func GenerateRandomBytes(length int) []byte {
	bytesBuffer := bytes.NewBuffer(nil)

	for l := 0; l < length; {
		n, _ := bytesBuffer.WriteString(rand.Text())
		l += n
	}

	return bytesBuffer.Bytes()[:length]
}
