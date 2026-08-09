package lib

import "github.com/mats0319/secure_transfer/utils"

type Encryptor interface {
	Encrypt(originFile string, encFile string) *utils.Error
}

type Decryptor interface {
	Decrypt(encFile string, decFile string) *utils.Error
}
