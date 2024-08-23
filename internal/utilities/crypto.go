package utilities

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

type SecureTokenOption interface{}

type withLength struct {
	length int
}

// WithLength Specifics length for SecureToken. If WithLength is not provided, default value is 16
func WithLength(length int) SecureTokenOption {
	return withLength{length: length}
}

// SecureToken creates a new random token
func SecureToken(options ...SecureTokenOption) string {
	length := 16
	for i := range options {
		option := options[i]

		if castedOpt, ok := option.(withLength); ok {
			length = castedOpt.length
		}
	}

	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
