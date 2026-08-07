package lib

import (
	"bufio"
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
}

func NewFileHeader(encMethod utils.EncryptMethod) *FileHeader {
	return &FileHeader{
		EncMethod: encMethod,
		Salt:      utils.GenerateRandomBytes[[]byte](utils.DeriveKeySaltLength),
		BaseNonce: utils.GenerateRandomBytes[[]byte](utils.AESBaseNonceLength),
		AAD:       []byte(utils.FileHeaderAAD),
	}
}

func (fh *FileHeader) Serialize() (fhBytes []byte, e *utils.Error) {
	e = fh.canSerialize()
	if e != nil {
		return
	}

	nonceStart := 7 + utils.DeriveKeySaltLength
	pubKStart := nonceStart + utils.AESBaseNonceLength
	aadStart := pubKStart + utils.PublicKeyLength
	length := aadStart + utils.FileHeaderAADLength

	fhBytes = make([]byte, length)

	copy(fhBytes[:2], utils.FileHeaderStart)
	fhBytes[2] = fh.EncMethod
	fhBytes[3] = byte(utils.DeriveKeySaltLength)
	fhBytes[4] = byte(utils.AESBaseNonceLength)
	fhBytes[5] = byte(utils.PublicKeyLength)
	fhBytes[6] = byte(utils.FileHeaderAADLength)

	copy(fhBytes[7:nonceStart], fh.Salt)
	copy(fhBytes[nonceStart:pubKStart], fh.BaseNonce)
	copy(fhBytes[pubKStart:aadStart], fh.PublicKey)
	copy(fhBytes[aadStart:], fh.AAD)

	return
}

func (fh *FileHeader) Deserialize(encFile string) (e *utils.Error) {
	file, err := os.Open(encFile)
	if err != nil {
		e = utils.ErrOpenFile().WithCause(err)
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
	if n != int(fileHeaderLen[0]) {
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
	case len(fh.AAD) != utils.FileHeaderAADLength:
		e = utils.ErrFileHeader().WithParam("aad", fh.AAD)
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

	nonceStart := 7 + utils.DeriveKeySaltLength
	aadStart := nonceStart + utils.AESBaseNonceLength
	pubKStart := aadStart + utils.FileHeaderAADLength
	length := pubKStart + utils.PublicKeyLength

	if len(data) != length {
		e = utils.ErrFileHeader().WithParam("length", len(data))
	}

	if e != nil {
		mlog.Error(e.String())
	}

	return
}
