package generate_avatar

import (
	"image"
	"image/color"
)

const basePX = 12

type imageImpl struct {
	displayColor    [5][3]bool
	color           color.RGBA
	backgroundColor color.RGBA
	clarity         int
}

func (m *imageImpl) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		return c
	})
}

func (m *imageImpl) Bounds() image.Rectangle {
	return image.Rect(0, 0, basePX*m.clarity, basePX*m.clarity)
}

func (m *imageImpl) At(x, y int) color.Color {
	// border
	if (0 <= x && x <= m.clarity-1) || (basePX*m.clarity-m.clarity-1 < x && x <= basePX*m.clarity-1) ||
		(0 <= y && y <= m.clarity-1) || (basePX*m.clarity-m.clarity-1 < y && y <= basePX*m.clarity-1) {
		return m.backgroundColor
	}

	// flip
	if x >= 7*m.clarity {
		return m.At(basePX*m.clarity-1-x, y)
	}

	// y - row, x - col
	// 矩阵上点的表示法与平面直角坐标系中点的表示法：
	//   [{r1, c1}, {r1, c2}]        [{x=1, y=3}, {x=2, y=3}]
	//   [{r2, c1}, {r2, c2}]   ->   [{x=1, y=2}, {x=2, y=2}]
	//   [{r3, c1}, {r3, c2}]        [{x=1, y=1}, {x=2, y=1}]
	// 举个例子，矩阵中第一列的点（col相同），在坐标系中的横坐标（x值）相同，所以x - col
	row := m.calcBlockPosition(y)
	col := m.calcBlockPosition(x)
	if !m.displayColor[row][col] {
		return m.backgroundColor
	}

	return m.color
}

// calcBlockPosition 计算色块位置，返回值不超过展示矩阵范围
func (m *imageImpl) calcBlockPosition(v int) int {
	// offset, skip border
	v -= m.clarity - 1

	index := 0
	for !(index*2*m.clarity < v && v <= (index+1)*2*m.clarity) {
		index++
	}

	return index
}
