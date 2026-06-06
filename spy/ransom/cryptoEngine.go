package ransom

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

type CryptoEngine struct {
	key []byte
}

func NewCryptoEngine(rawKey string) *CryptoEngine {
	// Хэшируем ключ через SHA-256, чтобы всегда получать стабильные 32 байта для AES-256
	hash := sha256.Sum256([]byte(rawKey))
	return &CryptoEngine{key: hash[:]}
}

func (ce *CryptoEngine) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ce.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Упаковываем: nonce + ciphertext
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (ce *CryptoEngine) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ce.key)
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
