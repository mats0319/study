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

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Go(func() { enc.read(ctx, originFrameCh, originFile) })
	wg.Go(func() {
		var wg2 sync.WaitGroup
		for range 8 {
			wg2.Go(func() { enc.encrypt(ctx, originFrameCh, encryptedFrameCh) })
		}
		wg2.Wait()

		close(encryptedFrameCh)
	})
	wg.Go(func() {
		e = enc.write(ctx, encryptedFrameCh, encFile)
		if e != nil {
			cancel()
		}
	})

	wg.Wait()

	return
}

func (enc *encryptorStream) init() (e *utils.Error) {
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

	// 检查：每一次加密文件都生成新的nonce，避免使用一个实例反复加密带来的安全风险
	enc.fileHeader, e = NewFileHeader(utils.EncryptMethod_Stream, tempPrivKey.PublicKey().Bytes())
	if e != nil {
		return
	}

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

func (enc *encryptorStream) read(ctx context.Context, ch chan *frame, originFile string) {
	defer close(ch)

	file, err := os.Open(originFile)
	if err != nil {
		e := utils.ErrOpenOriginFile().WithCause(err)
		mlog.Error(e.String())
		sendChanUnblocked(ctx, ch, &frame{e: e})
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, utils.FrameSize*2)
	counter := int32(0)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			originFrame := make([]byte, utils.FrameSize, utils.FrameSize+16)

			n, e := readExact(reader, originFrame)
			if e != nil {
				sendChanUnblocked(ctx, ch, &frame{e: e})
				return
			}

			isLastFrame := n < len(originFrame)

			sendChanUnblocked(ctx, ch, &frame{
				index:       counter,
				isLastFrame: isLastFrame,
				data:        originFrame[:n],
			})

			if isLastFrame {
				return
			}

			counter++
		}
	}
}

func (enc *encryptorStream) encrypt(ctx context.Context, originFrameCh chan *frame, encryptedFrameCh chan *frame) {
	frameIns := &frame{}
	ok := false

	for {
		select {
		case <-ctx.Done():
			return
		case frameIns, ok = <-originFrameCh:
			if !ok {
				return
			}
			if frameIns.e != nil {
				sendChanUnblocked(ctx, encryptedFrameCh, frameIns)
				return
			}

			nonce := makeNonce(enc.fileHeader.BaseNonce, frameIns.isLastFrame, frameIns.index)

			frameIns.data = enc.aesGCM.Seal(nil, nonce, frameIns.data, enc.fileHeader.AAD)

			sendChanUnblocked(ctx, encryptedFrameCh, frameIns)
		}
	}
}

func (enc *encryptorStream) write(ctx context.Context, ch chan *frame, encFile string) (e *utils.Error) {
	file, err := os.OpenFile(encFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		e = utils.ErrOpenEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	err = writer.WriteByte(byte(len(enc.fileHeader.Encoded)))
	if err != nil {
		e = utils.ErrWriteEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	_, err = writer.Write(enc.fileHeader.Encoded)
	if err != nil {
		e = utils.ErrWriteEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	frameCache := make(map[int32]*frame)
	writeIndex := int32(0)

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
	}

	err = writer.Flush()
	if err != nil {
		e = utils.ErrWriteEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}

func sendChanUnblocked(ctx context.Context, ch chan *frame, frameIns *frame) {
	select {
	case <-ctx.Done():
		return
	case ch <- frameIns:
		// avoid send block even deadlock, when channel buffer is full
	}
}
