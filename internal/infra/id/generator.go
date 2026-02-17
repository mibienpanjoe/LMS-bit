package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

func (Generator) NewID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("id-%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(b)
}
