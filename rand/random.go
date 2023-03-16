package rand

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func randBytes(l int) ([]byte, error) {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func String(length int) (string, error) {
	rb, err := randBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rb), nil
}

func Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// base64 encode the data into a string
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
