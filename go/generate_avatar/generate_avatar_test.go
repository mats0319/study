package generate_avatar

import (
	"math/rand/v2"
	"os"
	"testing"
)

func TestGenerateAvatar(t *testing.T) {
	testCase := []string{
		"mario",
		"mats0319",
		"avatar",
	}

	// 'rm -rf' out dir
	if err := os.RemoveAll("./images/"); err != nil {
		t.Error("empty out dir failed, error: ", err)
	}

	if err := os.MkdirAll("./images/", 0777); err != nil {
		t.Error("'mkdir' on out dir failed, error: ", err)
	}

	for i := range testCase {
		err := GenerateAvatar(testCase[i], rand.IntN(5)+3) // [3,7]
		if err != nil {
			t.Error("generate avatar failed, error:", err)
		}
	}
}

func TestDisplayColorMethod(t *testing.T) {
	str := "0123456789aAbBcCdDeEfF" // 16进制字符串通常使用小写字符，这里验证发现大小写不影响最后一个bit，不需要强制使用大写/小写字母

	box1 := make([]string, 0)
	box2 := make([]string, 0)
	for i := range str {
		isOdd := str[i]&0b01 == 1

		t.Logf("char: %c, byte: %b, is odd: %t", str[i], str[i], isOdd)

		if isOdd {
			box1 = append(box1, string(str[i]))
		} else {
			box2 = append(box2, string(str[i]))
		}
	}

	// 结论：根据最后一个bit是0或1，决定是否显示颜色的规则，显示颜色的概率为50%，
	// 即在足够的样本数量下，色块和白块的数量比为1：1
	t.Logf("box1: %v, box2: %v, box length: %d, %d", box1, box2, len(box1), len(box2))

	//    generate_avatar_test.go:41: char: 0, byte: 110000, is odd: false
	//    generate_avatar_test.go:41: char: 1, byte: 110001, is odd: true
	//    generate_avatar_test.go:41: char: 2, byte: 110010, is odd: false
	//    generate_avatar_test.go:41: char: 3, byte: 110011, is odd: true
	//    generate_avatar_test.go:41: char: 4, byte: 110100, is odd: false
	//    generate_avatar_test.go:41: char: 5, byte: 110101, is odd: true
	//    generate_avatar_test.go:41: char: 6, byte: 110110, is odd: false
	//    generate_avatar_test.go:41: char: 7, byte: 110111, is odd: true
	//    generate_avatar_test.go:41: char: 8, byte: 111000, is odd: false
	//    generate_avatar_test.go:41: char: 9, byte: 111001, is odd: true
	//    generate_avatar_test.go:41: char: a, byte: 1100001, is odd: true
	//    generate_avatar_test.go:41: char: A, byte: 1000001, is odd: true
	//    generate_avatar_test.go:41: char: b, byte: 1100010, is odd: false
	//    generate_avatar_test.go:41: char: B, byte: 1000010, is odd: false
	//    generate_avatar_test.go:41: char: c, byte: 1100011, is odd: true
	//    generate_avatar_test.go:41: char: C, byte: 1000011, is odd: true
	//    generate_avatar_test.go:41: char: d, byte: 1100100, is odd: false
	//    generate_avatar_test.go:41: char: D, byte: 1000100, is odd: false
	//    generate_avatar_test.go:41: char: e, byte: 1100101, is odd: true
	//    generate_avatar_test.go:41: char: E, byte: 1000101, is odd: true
	//    generate_avatar_test.go:41: char: f, byte: 1100110, is odd: false
	//    generate_avatar_test.go:41: char: F, byte: 1000110, is odd: false
	//    generate_avatar_test.go:52: box1: [1 3 5 7 9 a A c C e E], box2: [0 2 4 6 8 b B d D f F], box length: 11, 11
}
