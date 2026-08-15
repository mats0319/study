package generate_avatar

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
)

const basePX = 12
const minOpacity = 25 // opacity: 0.1
const lowOpacity = 64 // opacity <= 'min opacity', set it to 'low opacity'

type ImageImpl struct {
	DisplayColorFlag [5][3]bool
	Color            color.RGBA
	BackgroundColor  color.RGBA

	HashBytes []byte // hash bytes, 32 length
	Size      int
}

var _ image.Image = (*ImageImpl)(nil)

func NewImageImpl(text string, size int) *ImageImpl {
	if !(0 < size && size <= 10) {
		size = 1
	}

	ins := &ImageImpl{
		BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		HashBytes:       calcSHA256(text), // sha256结果必然是32位，后续代码可以依赖该值
		Size:            size,
	}

	{
		displayColorFlag := [5][3]bool{}
		for i := range displayColorFlag {
			for j := range displayColorFlag[i] {
				displayColorFlag[i][j] = ins.HashBytes[3*i+j]&0b01 == 1
			}
		}

		ins.DisplayColorFlag = displayColorFlag
	}

	{
		rgba := color.RGBA{
			R: ins.HashBytes[15],
			G: ins.HashBytes[16],
			B: ins.HashBytes[17],
			A: ins.HashBytes[18],
		}
		if rgba.A < minOpacity {
			rgba.A = lowOpacity
		}

		ins.Color = rgba
	}

	return ins
}

func (i *ImageImpl) FileName() string {
	return fmt.Sprintf("%x_%d.png", i.HashBytes[:6], i.Size)
}

func calcSHA256(v string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(v))

	return hasher.Sum(nil)
}

func (i *ImageImpl) ColorModel() color.Model {
	return color.RGBAModel
}

func (i *ImageImpl) Bounds() image.Rectangle {
	return image.Rect(0, 0, basePX*i.Size, basePX*i.Size)
}

// At x, y 是从0开始的索引
func (i *ImageImpl) At(x, y int) color.Color {
	// border
	if x < i.Size || x >= (5*2+1)*i.Size || y < i.Size || y >= (5*2+1)*i.Size {
		return i.BackgroundColor
	}

	// y - row, x - col
	// 矩阵上的点与平面直角坐标系中的点的对应关系：
	//   [{r1, c1}, {r1, c2}]        [{x=1, y=3}, {x=2, y=3}]
	//   [{r2, c1}, {r2, c2}]   ->   [{x=1, y=2}, {x=2, y=2}]
	//   [{r3, c1}, {r3, c2}]        [{x=1, y=1}, {x=2, y=1}]
	// 举个例子，矩阵中第一列的点（col相同），在坐标系中的横坐标（x值）相同，所以x - col
	// 这里其实y还应该做一次上下颠倒，然而我们认为这对生成图片的影响是可以接受的（相当于图片沿x轴翻转），所以不处理

	// 计算一个横/纵坐标值'v'在'size'等级下属于第几个块（结果为块索引）
	// 因为图片整体是正方形，所以横纵坐标可以使用相同规则处理
	row := (y - i.Size) / (2 * i.Size)
	col := (x - i.Size) / (2 * i.Size)
	col = min(col, 4-col) // flip right side

	var c color.Color
	if i.DisplayColorFlag[row][col] {
		c = i.Color
	} else {
		c = i.BackgroundColor
	}

	return c
}
