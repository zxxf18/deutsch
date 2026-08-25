package passwordcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const prefix = "v1"

type Cipher struct {
	aead cipher.AEAD
}

func New(key string) (*Cipher, error) {
	if key == "" {
		return nil, errors.New("password encryption key is required")
	}
	derivedKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(derivedKey[:])
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create AES-GCM: %w", err)
	}
	return &Cipher{aead: aead}, nil
}

func (c *Cipher) Encrypt(plaintext string) (string, error) {
	if c == nil || c.aead == nil {
		return "", errors.New("password cipher is not initialized")
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	sealed := c.aead.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, sealed...)
	return prefix + ":" + base64.RawStdEncoding.EncodeToString(payload), nil
}

func (c *Cipher) Decrypt(encoded string) (string, error) {
	if c == nil || c.aead == nil {
		return "", errors.New("password cipher is not initialized")
	}
	version, data, ok := strings.Cut(encoded, ":")
	if !ok || version != prefix {
		return "", errors.New("unsupported encrypted password format")
	}
	payload, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New("invalid encrypted password encoding")
	}
	if len(payload) < c.aead.NonceSize()+c.aead.Overhead() {
		return "", errors.New("encrypted password payload is too short")
	}
	nonce, ciphertext := payload[:c.aead.NonceSize()], payload[c.aead.NonceSize():]
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decrypt password: authentication failed")
	}
	return string(plaintext), nil
}

func (c *Cipher) Matches(encoded, candidate string) bool {
	plaintext, err := c.Decrypt(encoded)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(plaintext), []byte(candidate)) == 1
}
