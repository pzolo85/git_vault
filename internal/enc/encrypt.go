package enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func Encrypt(data, password []byte) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive a key from the password using Argon2id
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Generate a random nonce for AES-GCM
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	// Create AES-GCM instance
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create AES-GCM: %w", err)
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	// Combine salt, nonce, and ciphertext into a single slice
	finalData := append(salt, nonce...)
	finalData = append(finalData, ciphertext...)

	// Return the result as a base64-encoded string
	return base64.StdEncoding.EncodeToString(finalData), nil
}

// Decrypt decrypts the base64-encoded encrypted data using the password and returns the original data.
func Decrypt(encryptedData string, password string) ([]byte, error) {
	// Decode the base64-encoded data
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Extract the salt, nonce, and ciphertext
	if len(data) < 16+12 {
		return nil, errors.New("invalid encrypted data")
	}
	salt := data[:16]
	nonce := data[16:28]
	ciphertext := data[28:]

	// Derive the key from the password using Argon2id
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	// Create AES-GCM instance
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES-GCM: %w", err)
	}

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}
