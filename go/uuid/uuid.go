package uuid

import (
	"crypto/sha256"
	"github.com/google/uuid"
)

// New return uuid v4 string,
// with same 'data', it will return same string,
// without 'data', it will return random string.
func New(data ...byte) string {
	if len(data) < 1 {
		return uuid.NewString()
	}

	return uuid.NewHash(sha256.New(), uuid.Nil, data, 4).String()
}
