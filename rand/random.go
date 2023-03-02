package rand

import (
	"crypto/rand"
	"encoding/base64"
)

func hash(l int) ([]byte, error) {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func String(length int) (string, error) {
	rb, err := hash(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rb), nil
}
