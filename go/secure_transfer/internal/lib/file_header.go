package lib

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type FileHeader struct {
	EncMethod utils.EncryptMethod
	Salt      []byte // for derive key
	BaseNonce []byte // aes-gcm base nonce
	PublicKey []byte
	AAD       []byte

	Encoded []byte
}

func NewFileHeader(encMethod utils.EncryptMethod, pubKey []byte) (*FileHeader, *utils.Error) {
	fh := &FileHeader{
		EncMethod: encMethod,
		Salt:      utils.GenerateRandomBytes(utils.DeriveKeySaltLength),
		BaseNonce: utils.GenerateRandomBytes(utils.AESBaseNonceLength),
		PublicKey: pubKey,
	}

	e := fh.serialize()
	if e != nil {
		return nil, e
	}

	return fh, nil
}

func (fh *FileHeader) serialize() (e *utils.Error) {
	e = fh.canSerialize()
	if e != nil {
		return
	}

	nonceStart := 7 + utils.DeriveKeySaltLength
	pubKStart := nonceStart + utils.AESBaseNonceLength
	aadStart := pubKStart + utils.PublicKeyLength
	length := aadStart + utils.FileHeaderAADLength

	fh.Encoded = make([]byte, length)

	copy(fh.Encoded[:2], utils.FileHeaderStart)
	fh.Encoded[2] = fh.EncMethod
	fh.Encoded[3] = byte(utils.DeriveKeySaltLength)
	fh.Encoded[4] = byte(utils.AESBaseNonceLength)
	fh.Encoded[5] = byte(utils.PublicKeyLength)
	fh.Encoded[6] = byte(utils.FileHeaderAADLength)

	copy(fh.Encoded[7:nonceStart], fh.Salt)
	copy(fh.Encoded[nonceStart:pubKStart], fh.BaseNonce)
	copy(fh.Encoded[pubKStart:aadStart], fh.PublicKey)

	fh.AAD = utils.CalcSHA256(fh.Encoded[:aadStart])
	copy(fh.Encoded[aadStart:], fh.AAD)

	return
}

func (fh *FileHeader) Deserialize(encFile string) (e *utils.Error) {
	file, err := os.Open(encFile)
	if err != nil {
		e = utils.ErrOpenEncFile().WithCause(err)
		mlog.Error(e.String())
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, utils.FrameSize*2)

	fileHeaderLen := make([]byte, 1)
	_, e = readExact(reader, fileHeaderLen)
	if e != nil {
		return e
	}

	data := make([]byte, fileHeaderLen[0])
	n, e := readExact(reader, data)
	if e != nil {
		return
	}
	if n != len(data) {
		e = utils.ErrFileHeader().
			WithParam("length", n).
			WithParam("want", fileHeaderLen[0])
		mlog.Error(e.String())
		return
	}

	e = fh.canDeserialize(data)
	if e != nil {
		return
	}

	nonceStart := 7 + utils.DeriveKeySaltLength
	pubKStart := nonceStart + utils.AESBaseNonceLength
	aadStart := pubKStart + utils.PublicKeyLength

	fh.EncMethod = data[2]
	fh.Salt = data[7:nonceStart]
	fh.BaseNonce = data[nonceStart:pubKStart]
	fh.PublicKey = data[pubKStart:aadStart]
	fh.AAD = data[aadStart:]

	if !utils.CompareBytes(fh.AAD, utils.CalcSHA256(data[:aadStart])) {
		e = utils.ErrFileHeader().WithParam("hash not matched", "")
		mlog.Error(e.String())
		return
	}

	return
}

func (fh *FileHeader) canSerialize() (e *utils.Error) {
	switch {
	case !utils.IsEncryptMethod(fh.EncMethod):
		e = utils.ErrFileHeader().WithParam("enc method", fh.EncMethod)
	case len(fh.Salt) != utils.DeriveKeySaltLength:
		e = utils.ErrFileHeader().WithParam("salt", fh.Salt)
	case len(fh.BaseNonce) != utils.AESBaseNonceLength:
		e = utils.ErrFileHeader().WithParam("base_nonce", fh.BaseNonce)
	case len(fh.PublicKey) != utils.PublicKeyLength:
		e = utils.ErrFileHeader().WithParam("public_key", fh.PublicKey)
	}

	if e != nil {
		mlog.Error(e.String())
	}

	return
}

func (fh *FileHeader) canDeserialize(data []byte) (e *utils.Error) {
	switch {
	case len(data) < 7:
		e = utils.ErrFileHeader().WithParam("fixed length", len(data))
	case string(data[:2]) != utils.FileHeaderStart:
		e = utils.ErrFileHeader().WithParam("start", string(data[:2]))
	case !utils.IsEncryptMethod(data[2]):
		e = utils.ErrFileHeader().WithParam("enc method", data[2])
	case data[3] != utils.DeriveKeySaltLength:
		e = utils.ErrFileHeader().WithParam("salt len", data[3])
	case data[4] != utils.AESBaseNonceLength:
		e = utils.ErrFileHeader().WithParam("base nonce len", data[4])
	case data[5] != utils.PublicKeyLength:
		e = utils.ErrFileHeader().WithParam("pubK len", data[5])
	case data[6] != uint8(utils.FileHeaderAADLength):
		e = utils.ErrFileHeader().WithParam("aad len", data[6])
	}

	if e != nil {
		mlog.Error(e.String())
	}

	return
}

func skipFileHeader(reader io.Reader) (n int64, e *utils.Error) {
	fileHeaderLen := make([]byte, 1)
	_, e = readExact(reader, fileHeaderLen)
	if e != nil {
		return
	}

	fileHeader := make([]byte, fileHeaderLen[0])
	_, e = readExact(reader, fileHeader)
	if e != nil {
		return
	}

	n = 1 + int64(fileHeaderLen[0])

	return
}

func readExact(reader io.Reader, buf []byte) (n int, e *utils.Error) {
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		switch {
		case err == io.EOF:
			// 不视为错误
			break
		case errors.Is(err, io.ErrUnexpectedEOF):
			e = utils.ErrReadExact().WithParam("File truncated", "文件被截断")
			mlog.Error(e.String())
		default:
			e = utils.ErrReadExact().WithCause(err)
			mlog.Error(e.String())
		}
	}

	return
}
