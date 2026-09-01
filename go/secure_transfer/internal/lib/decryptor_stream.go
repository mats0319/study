package lib

import (
	"bufio"
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Go(func() { dec.read(ctx, encryptedFrameCh, encFile) })
	wg.Go(func() {
		var wg2 sync.WaitGroup
		for range 8 {
			wg2.Go(func() { dec.decrypt(ctx, encryptedFrameCh, decryptedFrameCh) })
		}
		wg2.Wait()

		close(decryptedFrameCh)
	})
	wg.Go(func() {
		e = dec.write(ctx, decryptedFrameCh, decFile)
		if e != nil {
			cancel()
		}
	})

	wg.Wait()

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

func (dec *decryptorStream) read(ctx context.Context, ch chan *frame, encFile string) {
	defer close(ch)

	file, err := os.Open(encFile)
	if err != nil {
		e := utils.ErrOpenEncFile().WithCause(err)
		mlog.Error(e.String())
		sendChanUnblocked(ctx, ch, &frame{e: e})
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, utils.FrameSize*2)

	_, e := skipFileHeader(reader)
	if e != nil {
		sendChanUnblocked(ctx, ch, &frame{e: e})
		return
	}

	for index := 0; ; index++ {
		select {
		case <-ctx.Done():
			return
		default:
			frameIns := &frame{}

			e = frameIns.deserialize(reader, index)
			if e != nil {
				sendChanUnblocked(ctx, ch, &frame{e: e})
				return
			}

			sendChanUnblocked(ctx, ch, frameIns)

			if frameIns.isLastFrame {
				return
			}
		}
	}
}

func (dec *decryptorStream) decrypt(ctx context.Context, encryptedFrameCh chan *frame, decryptedFrameCh chan *frame) {
	frameIns := &frame{}
	ok := false

	for {
		select {
		case <-ctx.Done():
			return
		case frameIns, ok = <-encryptedFrameCh:
			if !ok {
				return
			}
			if frameIns.e != nil {
				sendChanUnblocked(ctx, decryptedFrameCh, frameIns)
				return
			}

			nonce := makeNonce(dec.fileHeader.BaseNonce, frameIns.isLastFrame, frameIns.index)

			decryptedBytes, err := dec.aesGCM.Open(nil, nonce, frameIns.data, dec.fileHeader.AAD)
			if err != nil {
				e := utils.ErrAESDecrypt().WithCause(err)
				mlog.Error(e.String())
				sendChanUnblocked(ctx, decryptedFrameCh, &frame{e: e})
				return
			}

			frameIns.data = decryptedBytes

			sendChanUnblocked(ctx, decryptedFrameCh, frameIns)
		}
	}
}

func (dec *decryptorStream) write(ctx context.Context, ch chan *frame, decFile string) (e *utils.Error) {
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

ALL:
	for {
		select {
		case <-ctx.Done():
			return
		case frameIns, ok := <-ch:
			if !ok {
				break ALL
			}
			if frameIns.e != nil {
				e = frameIns.e
				return
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
	}

	err = writer.Flush()
	if err != nil {
		e = utils.ErrWriteDecFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
