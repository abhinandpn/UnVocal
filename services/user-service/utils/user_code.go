package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateUserCode() (string, error) {
	const prefix = "USR"
	const length = 8

	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}

	return prefix + string(b), nil
}
