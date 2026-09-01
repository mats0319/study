package lib

import (
	"fmt"
	"io"

	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

type frame struct {
	index       int32
	isLastFrame bool
	data        []byte
	e           *utils.Error
}

func (f *frame) serialize() []byte {
	res := make([]byte, 9+len(f.data))
	if f.isLastFrame {
		res[0] = 1
	} else {
		res[0] = 0
	}
	res[1] = byte(f.index >> 24)
	res[2] = byte(f.index >> 16)
	res[3] = byte(f.index >> 8)
	res[4] = byte(f.index)

	length := int32(len(f.data))
	res[5] = byte(length >> 24)
	res[6] = byte(length >> 16)
	res[7] = byte(length >> 8)
	res[8] = byte(length)

	copy(res[9:], f.data)

	return res
}

func (f *frame) deserialize(reader io.Reader, index int) (e *utils.Error) {
	fixed := make([]byte, 9)
	n, e := readExact(reader, fixed)
	if e != nil {
		return
	}
	if n < 9 {
		e = utils.ErrFrame().WithParam("length", n).WithParam("want", 9)
		mlog.Error(e.String())
		return
	}

	switch fixed[0] {
	case 0:
		f.isLastFrame = false
	case 1:
		f.isLastFrame = true
	default:
		e = utils.ErrFrame().WithParam("last flag", fixed[0])
		mlog.Error(e.String())
		return
	}

	f.index = (int32(fixed[1]) << 24) + (int32(fixed[2]) << 16) + (int32(fixed[3]) << 8) + int32(fixed[4])
	if f.index != int32(index) {
		e = utils.ErrFrame().WithParam(fmt.Sprintf("invalid frame index, want: %d, get: %d", index, f.index), "")
		mlog.Error(e.String())
		return
	}

	length := (int32(fixed[5]) << 24) + (int32(fixed[6]) << 16) + (int32(fixed[7]) << 8) + int32(fixed[8])
	if length > utils.FrameSize+16 {
		e = utils.ErrFrame().WithParam("should less than 1M, get: ", length)
		mlog.Error(e.String())
		return
	}

	data := make([]byte, length)
	n, e = readExact(reader, data)
	if e != nil {
		return
	}

	f.data = data[:n]

	return
}

func makeNonce(baseNonce []byte, isLastFrame bool, counter int32) []byte {
	nonce := make([]byte, utils.AESBaseNonceLength+1+4)
	copy(nonce[:utils.AESBaseNonceLength], baseNonce)

	if isLastFrame {
		nonce[utils.AESBaseNonceLength] = 1
	} else {
		nonce[utils.AESBaseNonceLength] = 0
	}

	nonce[utils.AESBaseNonceLength+1] = byte(counter >> 24)
	nonce[utils.AESBaseNonceLength+2] = byte(counter >> 16)
	nonce[utils.AESBaseNonceLength+3] = byte(counter >> 8)
	nonce[utils.AESBaseNonceLength+4] = byte(counter)

	return nonce
}
