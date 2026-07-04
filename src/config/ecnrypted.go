package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"os"
)

type EncryptedStorage struct {
	path string
	key  [32]byte
}

func NewEncryptedStorage(path, key string) Storage {
	return &EncryptedStorage{
		path: path,
		key:  sha256.Sum256([]byte(key)),
	}
}

func (s *EncryptedStorage) Get() (string, error) {
	fileData, err := os.ReadFile(s.path)
	if err != nil {
		return "", err
	}
	decryptedData, err := decrypt(fileData, s.key)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}

func decrypt(ciphertext []byte, key [32]byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("некорректный размер зашифрованных данных")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, actualCiphertext, nil)
}
