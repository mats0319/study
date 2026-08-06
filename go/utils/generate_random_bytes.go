package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand/v2"
)

const CharactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示全部字符库中的字符
const useMask = 1<<useBits - 1

// GenerateRandomBytes_CharacterLibraryIndex generate random 'length' readable Bytes
func GenerateRandomBytes_CharacterLibraryIndex(length int) []byte {
	b := make([]byte, length)

	randomNum, remainBits := rand.Int64(), 64
	for i := 0; i < len(b); {
		if remainBits < useBits {
			randomNum, remainBits = rand.Int64(), 64
		}

		index := int(randomNum & useMask) // 0b0011 1111
		if index < len(CharactersLibrary) {
			randomNum >>= useBits
			remainBits -= useBits

			b[i] = CharactersLibrary[index]
			i++
		} else {
			randomNum >>= 1
			remainBits -= 1
		}
	}

	return b
}

func GenerateRandomBytes_BytesEncode(length int) string {
	b := make([]byte, (length+1)/2)
	_, _ = crand.Read(b)

	return hex.EncodeToString(b)[:length]
}
