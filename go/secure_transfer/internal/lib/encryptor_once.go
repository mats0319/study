package lib

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/sha256"
	"os"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type encryptorOnce struct {
	fileHeader *FileHeader
	aesGCM     cipher.AEAD
	nonce      []byte

	pubKey *ecdh.PublicKey
}

func NewEncryptorOnce(pubKey *ecdh.PublicKey) Encryptor {
	return &encryptorOnce{pubKey: pubKey}
}

func (enc *encryptorOnce) Encrypt(originFile string, encFile string) *utils.Error {
	e := enc.init()
	if e != nil {
		return e
	}

	originFileBytes, err := os.ReadFile(originFile)
	if err != nil {
		e := utils.ErrReadOriginFile().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	encryptedBytes := enc.aesGCM.Seal(nil, enc.nonce, originFileBytes, enc.fileHeader.AAD)

	fileHeaderBytes, e := enc.fileHeader.Serialize()
	if e != nil {
		return e
	}

	file, err := os.OpenFile(encFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		e := utils.ErrOpenEncFile().WithCause(err)
		mlog.Error(e.String())
		return e
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_ = writer.WriteByte(byte(len(fileHeaderBytes)))
	_, _ = writer.Write(fileHeaderBytes)
	_, _ = writer.Write(encryptedBytes)

	err = writer.Flush()
	if err != nil {
		e := utils.ErrWriteEncFile().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func (enc *encryptorOnce) init() (e *utils.Error) {
	// 检查：每一次加密文件都生成新的nonce，避免使用一个实例反复加密带来的安全风险
	enc.fileHeader = NewFileHeader(utils.EncryptMethod_Once)

	// ecdh
	tempPrivKey, err := utils.Curve().GenerateKey(nil)
	if err != nil {
		e = utils.ErrECDH().WithCause(err)
		mlog.Error(e.String())
		return
	}

	sharedKey, err := tempPrivKey.ECDH(enc.pubKey)
	if err != nil {
		e = utils.ErrECDH().WithCause(err)
		mlog.Error(e.String())
		return
	}

	enc.fileHeader.PublicKey = tempPrivKey.PublicKey().Bytes()

	// kdf
	aesKey, err := hkdf.Key(sha256.New, sharedKey, enc.fileHeader.Salt, utils.DeriveKeyInfo, utils.DeriveKeyLength)
	if err != nil {
		e = utils.ErrKDF().WithCause(err)
		mlog.Error(e.String())
		return
	}

	// generate aes-gcm encryptor
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		e = utils.ErrAESNewCipher().WithCause(err)
		mlog.Error(e.String())
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		e = utils.ErrAESNewGCM().WithCause(err)
		mlog.Error(e.String())
		return
	}

	enc.aesGCM = aesGCM

	// nonce
	nonce := make([]byte, aesGCM.NonceSize()) // default zero placeholder
	copy(nonce[:utils.AESBaseNonceLength], enc.fileHeader.BaseNonce)

	enc.nonce = nonce

	return
}
