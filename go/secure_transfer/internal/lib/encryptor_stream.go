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

type encryptorStream struct {
	fileHeader *FileHeader
	aesGCM     cipher.AEAD

	pubKey *ecdh.PublicKey
}

func NewEncryptorStream(pubKey *ecdh.PublicKey) Encryptor {
	return &encryptorStream{pubKey: pubKey}
}

func (enc *encryptorStream) Encrypt(originFile string, encFile string) (e *utils.Error) {
	e = enc.init()
	if e != nil {
		return
	}

	originFrameCh := make(chan *frame, 100)
	encryptedFrameCh := make(chan *frame, 100)
	errCh := make(chan *utils.Error, 10)

	var wg sync.WaitGroup
	wg.Go(func() {
		e := enc.read(originFrameCh, originFile)
		if e != nil {
			errCh <- e
		}
		close(originFrameCh)
	})
	wg.Go(func() {
		var wg2 sync.WaitGroup
		for range 8 {
			wg2.Go(func() {
				e := enc.encrypt(originFrameCh, encryptedFrameCh)
				if e != nil {
					errCh <- e
				}
			})
		}
		wg2.Wait()
		close(encryptedFrameCh)
	})
	wg.Go(func() {
		e := enc.write(encryptedFrameCh, encFile)
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

func (enc *encryptorStream) init() (e *utils.Error) {
	// 检查：每一次加密文件都生成新的nonce，避免使用一个实例反复加密带来的安全风险
	enc.fileHeader = NewFileHeader(utils.EncryptMethod_Stream)

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

	return
}

func (enc *encryptorStream) read(ch chan *frame, originFile string) (e *utils.Error) {
	file, err := os.Open(originFile)
	if err != nil {
		e = utils.ErrOpenOriginFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, utils.FrameSize*2)
	counter := int32(0)

	for {
		originFrame := make([]byte, utils.FrameSize, utils.FrameSize+16)

		n, e := readExact(reader, originFrame)
		if e != nil {
			return e
		}

		isLastFrame := n < len(originFrame)

		ch <- &frame{
			index:       counter,
			isLastFrame: isLastFrame,
			data:        originFrame[:n],
		}

		if isLastFrame {
			break
		}

		counter++
	}

	return
}

func (enc *encryptorStream) encrypt(originFrameCh chan *frame, encryptedFrameCh chan *frame) (e *utils.Error) {
	for {
		frameIns, ok := <-originFrameCh
		if !ok {
			break
		}

		nonce := makeNonce(enc.fileHeader.BaseNonce, frameIns.isLastFrame, frameIns.index)

		frameIns.data = enc.aesGCM.Seal(nil, nonce, frameIns.data, enc.fileHeader.AAD)

		encryptedFrameCh <- frameIns
	}

	return
}

func (enc *encryptorStream) write(ch chan *frame, encFile string) (e *utils.Error) {
	file, err := os.OpenFile(encFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		e = utils.ErrOpenEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	fileHeaderBytes, e := enc.fileHeader.Serialize()
	if e != nil {
		return
	}

	err = writer.WriteByte(byte(len(fileHeaderBytes)))
	if err != nil {
		e = utils.ErrWriteEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	_, err = writer.Write(fileHeaderBytes)
	if err != nil {
		e = utils.ErrWriteEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	frameCache := make(map[int32]*frame)
	writeIndex := int32(0)

	for {
		frameIns, ok := <-ch
		if !ok {
			break
		}

		frameCache[frameIns.index] = frameIns

		for frameIns != nil && frameIns.index == writeIndex {
			_, err := writer.Write(frameIns.serialize())
			if err != nil {
				e = utils.ErrWriteEncFile().WithCause(err)
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
		e = utils.ErrWriteEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
