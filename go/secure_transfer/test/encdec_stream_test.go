package test

import (
	"testing"

	"github.com/mats0319/secure_transfer/internal"
)

func TestEncDecStream(t *testing.T) {
	err := generateRandomFile("message.txt", 100)
	if err != nil {
		t.Error(err)
		return
	}

	err = internal.GenerateKeyPair()
	if err != nil {
		t.Error(err)
		return
	}

	err = internal.Encrypt()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("> Encrypt Success.")

	err = internal.Decrypt()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("> Decrypt Success.")

	isEqual, err := compareFile(t, "message.txt", "message_decrypted.txt")
	if err != nil || !isEqual {
		t.Error(err)
		return
	}

	t.Log("> File Hash Matched, Test Passed.")
}
