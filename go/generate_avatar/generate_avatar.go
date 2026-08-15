package generate_avatar

import (
	"bytes"
	"errors"
	"image/png"
)

// GenerateAvatar generate avatar according to input params
// - text: use for calculate hash, distinguish 'A' and 'a'
// - size: image size, valid range is [1, 10], image will be '12*size x 12*size' px
func GenerateAvatar(text string, size int) (fileName string, imageBytes []byte, err error) {
	if !(0 < size && size <= 10) {
		err = errors.New("invalid image size, required: (0, 10]")
		return
	}

	imageImplIns := NewImageImpl(text, size)

	// encode to png
	buffer := bytes.NewBuffer(nil)
	err = (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(buffer, imageImplIns)
	if err != nil {
		return
	}

	fileName = imageImplIns.FileName()
	imageBytes = buffer.Bytes()

	return
}
