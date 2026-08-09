package lib

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/sha256"
	"os"
	"sync"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type decryptorStream struct {
	aesGCM cipher.AEAD

	privKey    *ecdh.PrivateKey
	fileHeader *FileHeader
}

func NewDecryptorStream(privKey *ecdh.PrivateKey, fileHeader *FileHeader) Decryptor {
	return &decryptorStream{privKey: privKey, fileHeader: fileHeader}
}

func (dec *decryptorStream) Decrypt(encFile string, decFile string) (e *utils.Error) {
	e = dec.init()
	if e != nil {
		return
	}

	encryptedFrameCh := make(chan *frame, 100)
	decryptedFrameCh := make(chan *frame, 100)
	errCh := make(chan *utils.Error, 10)

	var wg sync.WaitGroup
	wg.Go(func() {
		e := dec.read(encryptedFrameCh, encFile)
		if e != nil {
			errCh <- e
		}
		close(encryptedFrameCh)
	})
	wg.Go(func() {
		var wg2 sync.WaitGroup
		for range 8 {
			wg2.Go(func() {
				e := dec.decrypt(encryptedFrameCh, decryptedFrameCh)
				if e != nil {
					errCh <- e
				}
			})
		}
		wg2.Wait()
		close(decryptedFrameCh)
	})
	wg.Go(func() {
		e := dec.write(decryptedFrameCh, decFile)
		if e != nil {
			errCh <- e
		}
	})

	wg.Wait()

	if len(errCh) > 0 {
		e = <-errCh
	}

	return
}

func (dec *decryptorStream) init() (e *utils.Error) {
	// ecdh
	pubKey, err := utils.Curve().NewPublicKey(dec.fileHeader.PublicKey)
	if err != nil {
		e = utils.ErrECDH().WithCause(err)
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

	return
}

func (dec *decryptorStream) read(ch chan *frame, encFile string) (e *utils.Error) {
	file, err := os.Open(encFile)
	if err != nil {
		e = utils.ErrOpenEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, utils.FrameSize*2)

	_, e = skipFileHeader(reader)
	if e != nil {
		return
	}

	for {
		frameIns := &frame{}

		e = frameIns.deserialize(reader)
		if e != nil {
			return
		}

		ch <- frameIns

		if frameIns.isLastFrame {
			break
		}
	}

	return
}

func (dec *decryptorStream) decrypt(encryptedFrameCh chan *frame, decryptedFrameCh chan *frame) (e *utils.Error) {
	for {
		frameIns, ok := <-encryptedFrameCh
		if !ok {
			break
		}

		nonce := makeNonce(dec.fileHeader.BaseNonce, frameIns.isLastFrame, frameIns.index)

		decryptedBytes, err := dec.aesGCM.Open(nil, nonce, frameIns.data, dec.fileHeader.AAD)
		if err != nil {
			e = utils.ErrAESDecrypt().WithCause(err)
			mlog.Error(e.String())
			return e
		}

		frameIns.data = decryptedBytes

		decryptedFrameCh <- frameIns
	}

	return
}

func (dec *decryptorStream) write(ch chan *frame, decFile string) (e *utils.Error) {
	frameCache := make(map[int32]*frame)
	writeIndex := int32(0)

	file, err := os.OpenFile(decFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		e = utils.ErrOpenDecFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for {
		frameIns, ok := <-ch
		if !ok {
			break
		}

		frameCache[frameIns.index] = frameIns

		for frameIns != nil && frameIns.index == writeIndex {
			_, err := writer.Write(frameIns.data)
			if err != nil {
				e = utils.ErrWriteDecFile().WithCause(err)
				mlog.Error(e.String())
				return
			}

			delete(frameCache, writeIndex)
			writeIndex++

			frameIns, _ = frameCache[writeIndex]
		}
	}

	err = writer.Flush()
	if err != nil {
		e = utils.ErrWriteDecFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
