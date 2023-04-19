package utils

import (
	"crypto/rand"
	"os"
)

func GenerateAndSaveEncryptionKey(filename string) ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filename, key, 0600)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func LoadEncryptionKey(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func ExistsEncryptionKey(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}
