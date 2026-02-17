package id

import (
	"crypto/rand"
	"encoding/hex"
)

type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

func (Generator) NewID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "id-fallback"
	}

	return hex.EncodeToString(b)
}
