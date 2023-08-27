package uuid

import (
	"fmt"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	for range [5]struct{}{} {
		fmt.Println(New())
	}
}

func TestGenerateUUIDByDesignatedText(t *testing.T) {
	text := []string{
		"config",
		"db",
	}

	expectedUID := []string{
		"612fbb57-c44c-4b54-8188-13d1dd598306",
		"8f7b02be-c349-4aee-8305-503580bbbf12",
	}

	if len(text) != len(expectedUID) {
		t.Log("unmatched params amount")
		t.Fail()
	}

	for i := range text {
		uid := New([]byte(text[i])...)
		if uid != expectedUID[i] {
			t.Logf("%s uid generate failed.\nget     : %s\nexpected: %s\n", text[i], uid, expectedUID)
			t.Fail()
		}
	}
}
