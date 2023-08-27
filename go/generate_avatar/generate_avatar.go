package generate_avatar

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

// GenerateAvatar generate avatar according to 'text'
//
//   @params text: use for calculate hash, distinguish 'A' and 'a'
//   @params clarity: image clarity level, valid range is (0, 10]
func GenerateAvatar(text string, clarity int) error {
	if !(0 < clarity && clarity <= 10) {
		return errors.New("invalid clarity, required: (0, 10]")
	}

	hash := calcSHA256(text)
	if len(hash) != 64 { // 其实这里可以不验长度，但是下面用到了长度数值，所以顺便验一下
		return errors.New(fmt.Sprintf("calc hash failed, expect 64 digits, get %d digits\n", len(hash)))
	}

	// calc image style
	imageImplIns := &imageImpl{
		clarity:         clarity,
		backgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}
	{
		matrix := [5][3]bool{}
		for i := 0; i < len(matrix); i++ {
			for j := 0; j < len(matrix[i]); j++ {
				matrix[i][j] = isOdd(hash[12*i+4*j+3])
			}
		}

		rgba, err := parseColor(hash[60:])
		if err != nil {
			return err
		}

		imageImplIns.displayColor = matrix
		imageImplIns.color = rgba
	}

	// encode to image
	writer := bytes.NewBufferString("")
	{
		b64 := base64.NewEncoder(base64.StdEncoding, writer)
		err := (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(b64, imageImplIns)
		if err != nil {
			return err
		}
		_ = b64.Close()
	}

	// serialize and write file
	byteSlice, err := base64.StdEncoding.DecodeString(writer.String())
	if err != nil {
		return err
	}

	err = os.MkdirAll("./img/", 0755)
	if err != nil {
		return err
	}

	fileName := text
	if len(fileName) > 8 {
		fileName = fileName[:8]
	}
	fileName = fmt.Sprintf("./img/%s_%d.png", fileName, clarity)

	return os.WriteFile(fileName, byteSlice, 0755)
}

// calcSHA256 always return a 64-length hex string
func calcSHA256(text string) string {
	hash := sha256.New()
	hash.Reset()
	hash.Write([]byte(text))
	byteSlice := hash.Sum(nil)

	return hex.EncodeToString(byteSlice)
}

func isOdd(char byte) bool {
	return char == '1' || char == '3' || char == '5' || char == '7' || char == '9' ||
		char == 'b' || char == 'B' || char == 'd' || char == 'D' || char == 'f' || char == 'F'
}

func parseColor(str string) (color.RGBA, error) {
	r, err := strconv.ParseInt(str[:2], 16, 0)
	g, err2 := strconv.ParseInt(str[1:3], 16, 0)
	b, err3 := strconv.ParseInt(str[2:], 16, 0)
	a, err4 := strconv.ParseInt(string([]byte{str[3], str[0]}), 16, 0)
	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		return color.RGBA{}, errors.New("parse color failed")
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
}
