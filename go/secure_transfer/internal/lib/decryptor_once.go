package lib

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/sha256"
	"io"
	"os"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type decryptorOnce struct {
	aesGCM cipher.AEAD
	nonce  []byte

	privKey     *ecdh.PrivateKey
	fileHeader  *FileHeader
	encFileSize int64
}

func NewDecryptorOnce(privKey *ecdh.PrivateKey, fileHeader *FileHeader, encFileSize int64) Decryptor {
	return &decryptorOnce{
		privKey:     privKey,
		fileHeader:  fileHeader,
		encFileSize: encFileSize,
	}
}

func (dec *decryptorOnce) Decrypt(encFile string, decFile string) (e *utils.Error) {
	e = dec.init()
	if e != nil {
		return
	}

	ciphertext, e := dec.getCiphertext(encFile)
	if e != nil {
		return
	}

	decryptedBytes, err := dec.aesGCM.Open(nil, dec.nonce, ciphertext, dec.fileHeader.AAD)
	if err != nil {
		e = utils.ErrAESDecrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	err = os.WriteFile(decFile, decryptedBytes, 0644)
	if err != nil {
		e = utils.ErrWriteFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}

func (dec *decryptorOnce) init() (e *utils.Error) {
	// ecdh
	pubKey, err := utils.Curve().NewPublicKey(dec.fileHeader.PublicKey)
	if err != nil {
		e = utils.ErrInvalidPublicKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	sharedKey, err := dec.privKey.ECDH(pubKey)
	if err != nil {
		e = utils.ErrECDH().WithCause(err)
		mlog.Error(e.String())
		return
	}

	// kdf
	aesKey, err := hkdf.Key(sha256.New, sharedKey, dec.fileHeader.Salt, utils.DeriveKeyInfo, utils.DeriveKeyLength)
	if err != nil {
		e = utils.ErrKDF().WithCause(err)
		mlog.Error(e.String())
		return
	}

	// generate aes-gcm decryptor
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

	dec.aesGCM = aesGCM

	// nonce
	nonce := make([]byte, aesGCM.NonceSize())
	copy(nonce[:utils.AESBaseNonceLength], dec.fileHeader.BaseNonce)

	dec.nonce = nonce

	return
}

func (dec *decryptorOnce) getCiphertext(encFile string) (ciphertext []byte, e *utils.Error) {
	file, err := os.Open(encFile)
	if err != nil {
		e = utils.ErrOpenFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	n, e := skipFileHeader(reader)
	if e != nil {
		return
	}
	if dec.encFileSize < n {
		e = utils.ErrEncryptedFile().
			WithParam("file length", dec.encFileSize).
			WithParam("header length", n)
		mlog.Error(e.String())
		return
	}

	ciphertext = make([]byte, dec.encFileSize-n)
	_, err = reader.Read(ciphertext)
	if err != nil {
		e = utils.ErrReadFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}

func skipFileHeader(reader io.Reader) (n int64, e *utils.Error) {
	fileHeaderLen := make([]byte, 1)
	_, err := reader.Read(fileHeaderLen)
	if err != nil {
		e = utils.ErrReadFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	fileHeader := make([]byte, fileHeaderLen[0])
	_, err = reader.Read(fileHeader)
	if err != nil {
		e = utils.ErrReadFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	n = 1 + int64(fileHeaderLen[0])

	return
}
