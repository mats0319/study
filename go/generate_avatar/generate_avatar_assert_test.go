package generate_avatar

// 本文件集中验证 generate_avatar 的核心性质（断言），全部在内存中完成，不产生任何文件写入。
// 覆盖内容：图片尺寸、PNG 有效性、白边、左右对称、色块一致性、flag 与渲染一致性、
// 输出确定性、文件名格式、无效尺寸处理、透明度下限、同文本不同尺寸、空文本。
// 运行方式: go test -v ./generate_avatar/

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"regexp"
	"testing"
)

var white = color.RGBA{R: 255, G: 255, B: 255, A: 255}

// decodePNG 解码 GenerateAvatar 输出的字节流，失败即终止当前测试
func decodePNG(t *testing.T, data []byte) image.Image {
	t.Helper()
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode png failed: %v", err)
	}
	return img
}

// sameColor 通过 16bit RGBA 分量比较两个颜色（可正确处理带透明度的颜色）
func sameColor(a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// --- 1. 尺寸：图片必须为 12*size x 12*size ---
func TestGeneratedImageSize(t *testing.T) {
	for _, size := range []int{1, 3, 10} { // 边界 + 常规
		_, data, err := GenerateAvatar("mario", size)
		if err != nil {
			t.Fatalf("size=%d: generate failed: %v", size, err)
		}
		img := decodePNG(t, data)
		if w, h := img.Bounds().Dx(), img.Bounds().Dy(); w != 12*size || h != 12*size {
			t.Errorf("size=%d: got %dx%d, want %dx%d", size, w, h, 12*size, 12*size)
		}
	}
}

// --- 2. PNG 有效性：字节流可解码且带合法 PNG 签名 ---
func TestGeneratedPNGIsValid(t *testing.T) {
	pngSig := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	for _, size := range []int{1, 4} {
		_, data, err := GenerateAvatar("mats0319", size)
		if err != nil {
			t.Fatalf("size=%d: generate failed: %v", size, err)
		}
		if !bytes.HasPrefix(data, pngSig) {
			t.Errorf("size=%d: bad png signature: %x", size, data[:min(len(data), 8)])
		}
		img := decodePNG(t, data)
		if img == nil {
			t.Fatalf("size=%d: decoded image is nil", size)
		}
	}
}

// --- 3. 白边：外圈宽度为 size 的区域必须是纯白 ---
func TestWhiteBorder(t *testing.T) {
	size := 3
	ins := NewImageImpl("mario", size)
	w, h := ins.Bounds().Dx(), ins.Bounds().Dy()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			inBorder := x < size || x >= 11*size || y < size || y >= 11*size
			if !inBorder {
				continue
			}
			if !sameColor(ins.At(x, y), white) {
				t.Fatalf("border pixel (%d,%d) is not white: %v", x, y, ins.At(x, y))
			}
		}
	}
}

// --- 4. 左右对称：任意像素与其水平镜像像素颜色一致 ---
func TestHorizontalSymmetry(t *testing.T) {
	for _, size := range []int{1, 4} {
		ins := NewImageImpl("mats0319", size)
		w, h := ins.Bounds().Dx(), ins.Bounds().Dy()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if !sameColor(ins.At(x, y), ins.At(w-1-x, y)) {
					t.Fatalf("size=%d: asymmetry at (%d,%d) vs (%d,%d)",
						size, x, y, w-1-x, y)
				}
			}
		}
	}
}

// --- 5. 色块一致性：每个 2*size 见方的格子内部颜色必须完全一致 ---
func TestBlockUniform(t *testing.T) {
	for _, size := range []int{1, 3} {
		ins := NewImageImpl("generate_avatar", size)
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				x0, y0 := size+c*2*size, size+r*2*size // 格子左上角
				base := ins.At(x0, y0)
				for dy := 0; dy < 2*size; dy++ {
					for dx := 0; dx < 2*size; dx++ {
						if !sameColor(base, ins.At(x0+dx, y0+dy)) {
							t.Fatalf("size=%d: block(%d,%d) not uniform at (%d,%d)",
								size, r, c, x0+dx, y0+dy)
						}
					}
				}
			}
		}
	}
}

// --- 6. flag 与渲染一致（白盒）：格子中心像素颜色必须由 DisplayColorFlag 决定 ---
func TestFlagMatchesRenderedColor(t *testing.T) {
	ins := NewImageImpl("mario", 3)
	for r := 0; r < 5; r++ {
		for c := 0; c < 3; c++ { // 左两列 + 中列；右两列由镜像保证
			x := size0Center(ins.Size, c)
			y := size0Center(ins.Size, r)
			want := color.Color(ins.BackgroundColor)
			if ins.DisplayColorFlag[r][c] {
				want = ins.Color
			}
			if !sameColor(ins.At(x, y), want) {
				t.Fatalf("flag[%d][%d]=%v but center pixel=%v", r, c, ins.DisplayColorFlag[r][c], ins.At(x, y))
			}
		}
	}
}

// size0Center 返回第 n 个格子（0 基）中心像素的坐标
func size0Center(size, n int) int {
	return size + n*2*size + size - 1
}

// --- 7. 确定性：相同输入产出完全相同的字节；不同输入产出不同字节 ---
func TestDeterministicOutput(t *testing.T) {
	_, d1, err := GenerateAvatar("mario", 4)
	if err != nil {
		t.Fatal(err)
	}
	_, d2, err := GenerateAvatar("mario", 4)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(d1, d2) {
		t.Fatal("same input must produce identical bytes")
	}

	// 不同文本（doc 注释强调区分大小写）
	_, dA, _ := GenerateAvatar("A", 3)
	_, da, _ := GenerateAvatar("a", 3)
	if bytes.Equal(dA, da) {
		t.Error("'A' and 'a' must produce different images")
	}
	_, dM, _ := GenerateAvatar("mario", 3)
	if bytes.Equal(dM, da) {
		t.Error("different texts must produce different images")
	}
}

// --- 8. 文件名格式：<12位hex>_<size>.png 且与哈希前6字节一致 ---
func TestFileNameFormat(t *testing.T) {
	re := regexp.MustCompile(`^[0-9a-f]{12}_\d+\.png$`)

	ins := NewImageImpl("mario", 3)
	want := fmt.Sprintf("%x_%d.png", ins.HashBytes[:6], ins.Size)
	if got := ins.FileName(); got != want {
		t.Fatalf("FileName() = %q, want %q", got, want)
	}
	if !re.MatchString(ins.FileName()) {
		t.Fatalf("FileName() = %q does not match %v", ins.FileName(), re)
	}

	// 同一文本不同尺寸必须产生不同文件名（避免互相覆盖）
	ins1 := NewImageImpl("mario", 1)
	ins3 := NewImageImpl("mario", 3)
	if ins1.FileName() == ins3.FileName() {
		t.Fatalf("different sizes must not share filename: %q", ins1.FileName())
	}
}

// --- 9. 无效尺寸：GenerateAvatar 拒绝；NewImageImpl 回退到默认尺寸 1 ---
func TestGenerateAvatarRejectsInvalidSize(t *testing.T) {
	for _, size := range []int{0, -1, 11} {
		if _, _, err := GenerateAvatar("mario", size); err == nil {
			t.Errorf("size=%d: expected error, got nil", size)
		}
	}
}

func TestNewImageImplDefaultsInvalidSize(t *testing.T) {
	for _, size := range []int{0, -5, 11} {
		ins := NewImageImpl("mario", size)
		if ins.Size != 1 {
			t.Errorf("size=%d: Size = %d, want default 1", size, ins.Size)
		}
		if w, h := ins.Bounds().Dx(), ins.Bounds().Dy(); w != 12 || h != 12 {
			t.Errorf("size=%d: bounds %dx%d, want 12x12", size, w, h)
		}
	}
}

// --- 10. 透明度下限：alpha < minOpacity 时置为 lowOpacity，最终 alpha 恒 >= minOpacity ---
func TestAlphaFloor(t *testing.T) {
	for _, text := range []string{"mario", "mats0319", "generate_avatar", "", "A", "a", "hello world"} {
		ins := NewImageImpl(text, 3)

		wantA := ins.HashBytes[18]
		if wantA < minOpacity {
			wantA = lowOpacity
		}
		if ins.Color.A != wantA {
			t.Errorf("text %q: alpha = %d, want %d", text, ins.Color.A, wantA)
		}
		if ins.Color.A < minOpacity {
			t.Errorf("text %q: alpha %d below minOpacity %d", text, ins.Color.A, minOpacity)
		}

		// 颜色 RGB 必须直接来自哈希的 15/16/17 字节
		if ins.Color.R != ins.HashBytes[15] || ins.Color.G != ins.HashBytes[16] || ins.Color.B != ins.HashBytes[17] {
			t.Errorf("text %q: color %v does not match hash bytes %x", text, ins.Color, ins.HashBytes[15:19])
		}
	}
}

// --- 11. 同文本不同尺寸：图案与颜色一致，仅尺寸（像素数）不同 ---
func TestSameTextDifferentSize(t *testing.T) {
	ins1 := NewImageImpl("mario", 1)
	ins3 := NewImageImpl("mario", 3)

	if ins1.DisplayColorFlag != ins3.DisplayColorFlag {
		t.Error("same text must produce the same pattern regardless of size")
	}
	if ins1.Color != ins3.Color {
		t.Errorf("same text must produce the same color, got %v vs %v", ins1.Color, ins3.Color)
	}
	if ins1.Bounds() == ins3.Bounds() {
		t.Error("different sizes must produce different bounds")
	}
}

// --- 12. 不同文本必须产生不同图案（防止图案恒定的回归） ---
func TestPatternsDifferBetweenTexts(t *testing.T) {
	ins1 := NewImageImpl("mario", 3)
	ins2 := NewImageImpl("mats0319", 3)

	differ := false
	for i := range ins1.DisplayColorFlag {
		for j := range ins1.DisplayColorFlag[i] {
			if ins1.DisplayColorFlag[i][j] != ins2.DisplayColorFlag[i][j] {
				differ = true
			}
		}
	}
	if !differ {
		t.Error("mario and mats0319 must produce different patterns")
	}
}

// --- 13. 空文本：不应报错，且生成合法图片 ---
func TestEmptyText(t *testing.T) {
	fileName, data, err := GenerateAvatar("", 3)
	if err != nil {
		t.Fatal("empty text should not error, got:", err)
	}
	if fileName == "" {
		t.Fatal("empty text should still produce a filename")
	}
	img := decodePNG(t, data)
	if w, h := img.Bounds().Dx(), img.Bounds().Dy(); w != 36 || h != 36 {
		t.Fatalf("empty text: got %dx%d, want 36x36", w, h)
	}
}
