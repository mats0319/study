package generate_avatar

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
)

const basePX = 12

var _ image.Image = (*ImageImpl)(nil)

type ImageImpl struct {
	DisplayColorFlag [5][3]bool
	Color            color.RGBA
	BackgroundColor  color.RGBA
	Size             int
}

func NewImageImpl(text string, size int) (*ImageImpl, error) {
	ins := &ImageImpl{
		BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Size:            size,
	}

	hash := calcSHA256(text)
	if len(hash) != 64 { // 其实这里可以不验长度，但是下面用到了长度数值，所以顺便验一下
		return nil, errors.New(fmt.Sprintf("calc hash failed, expect 64 digits, get %d digits\n", len(hash)))
	}

	{
		displayColorFlag := [5][3]bool{}
		for i := range len(displayColorFlag) {
			for j := range len(displayColorFlag[i]) {
				displayColorFlag[i][j] = hash[12*i+4*j+3]&0b01 == 1
			}
		}

		ins.DisplayColorFlag = displayColorFlag
	}

	{
		rgba, err := parseColor(hash[60:])
		if err != nil {
			return nil, err
		}

		ins.Color = rgba
	}

	return ins, nil
}

func (i *ImageImpl) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color { return c })
}

func (i *ImageImpl) Bounds() image.Rectangle {
	return image.Rect(0, 0, basePX*i.Size, basePX*i.Size)
}

func (i *ImageImpl) At(x, y int) color.Color {
	// border
	if (0 <= x && x <= i.Size-1) || (basePX*i.Size-i.Size-1 < x && x <= basePX*i.Size-1) ||
		(0 <= y && y <= i.Size-1) || (basePX*i.Size-i.Size-1 < y && y <= basePX*i.Size-1) {
		return i.BackgroundColor
	}

	// flip
	if x >= 7*i.Size {
		return i.At(basePX*i.Size-1-x, y)
	}

	// y - row, x - col
	// 矩阵上点的表示法与平面直角坐标系中点的表示法：
	//   [{r1, c1}, {r1, c2}]        [{x=1, y=3}, {x=2, y=3}]
	//   [{r2, c1}, {r2, c2}]   ->   [{x=1, y=2}, {x=2, y=2}]
	//   [{r3, c1}, {r3, c2}]        [{x=1, y=1}, {x=2, y=1}]
	// 举个例子，矩阵中第一列的点（col相同），在坐标系中的横坐标（x值）相同，所以x - col
	row := calcBlockPosition(y, i.Size)
	col := calcBlockPosition(x, i.Size)
	if !i.DisplayColorFlag[row][col] {
		return i.BackgroundColor
	}

	return i.Color
}

// calcBlockPosition 计算色块位置，返回值不超过展示矩阵范围
func calcBlockPosition(v int, size int) int {
	// offset, skip border
	v -= size - 1

	index := 0
	for !(index*2*size < v && v <= (index+1)*2*size) {
		index++
	}

	return index
}

func parseColor(str string) (color.RGBA, error) {
	if len(str) != 4 {
		return color.RGBA{}, errors.New("invalid color str")
	}

	r, err := strconv.ParseInt(str[:2], 16, 0)
	g, err2 := strconv.ParseInt(str[1:3], 16, 0)
	b, err3 := strconv.ParseInt(str[2:], 16, 0)
	a, err4 := strconv.ParseInt(string([]byte{str[3], str[0]}), 16, 0)
	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		return color.RGBA{}, errors.New("parse color failed")
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
}
