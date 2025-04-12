package util

import "crypto/rand"

func RandomString(length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return string(b), nil
}
