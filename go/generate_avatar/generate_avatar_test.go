package generate_avatar

import "testing"

func TestGenerateAvatar(t *testing.T) {
	testCase := []string{
		"mario",
		"mats9693",
	}

	for i := range testCase {
		err := GenerateAvatar(testCase[i], 3)
		if err != nil {
			t.Error("generate avatar failed", err)
		}
	}
}
