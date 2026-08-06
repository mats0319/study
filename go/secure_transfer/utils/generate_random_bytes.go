package utils

import "math/rand/v2"

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
