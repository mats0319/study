package generate_avatar

import (
	"os"
	"testing"
)

func TestGenerateSingleAvatar(t *testing.T) {
	err := os.MkdirAll("./images/", 0755)
	if err != nil {
		t.Fatal("create dir failed, error:", err)
	}

	writeAvatar(t, "mario", 1)
}

func TestGenerateAvatar(t *testing.T) {
	testCase := []string{
		"mario",
		"mats0319",
		"generate_avatar",
	}

	err := os.MkdirAll("./images/", 0755)
	if err != nil {
		t.Fatal("create dir failed, error:", err)
	}

	for i := range testCase {
		writeAvatar(t, testCase[i], 3)
	}
}

func writeAvatar(t *testing.T, text string, size int) {
	t.Helper()
	fileName, fileBytes, err := GenerateAvatar(text, size)
	if err != nil {
		t.Fatal("generate avatar failed, error:", err)
	}

	err = os.WriteFile("./images/"+fileName, fileBytes, 0644)
	if err != nil {
		t.Fatal("write file failed, error:", err)
	}

	t.Logf("Generate File: './images/%s'\n", fileName)
}
