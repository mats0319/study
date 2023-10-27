package generate_avatar

import (
	"os"
	"testing"
)

func TestGenerateAvatar(t *testing.T) {
	testCase := []string{
		"mario",
		"mats9693",
	}

	// 'rm -rf' out dir
	err := os.RemoveAll("./img/")
	if err != nil {
		t.Error("empty out dir failed, error: ", err)
	}

	err = os.MkdirAll("./img/", 0777)
	if err != nil {
		t.Error("'mkdir' on out dir failed, error: ", err)
	}

	for i := range testCase {
		err = GenerateAvatar(testCase[i], 4)
		if err != nil {
			t.Error("generate avatar failed, error:", err)
		}
	}
}
