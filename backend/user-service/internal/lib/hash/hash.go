package hash

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

const (
	SaltLength = 16
)

// hashPasswordWithSaltAdvanced хеширует пароль с солью с настройками
func HashPasswordWithSalt(password string) (hash string, salt string, err error) {
	// Генерируем соль
	saltBytes := make([]byte, SaltLength)
	_, err = rand.Read(saltBytes)
	if err != nil {
		return "", "", err
	}
	salt = base64.StdEncoding.EncodeToString(saltBytes)

	hasher := sha512.New()

	// Хешируем пароль + соль
	combined := password + salt
	hasher.Write([]byte(combined))
	hashBytes := hasher.Sum(nil)
	hash = hex.EncodeToString(hashBytes)

	return hash, salt, nil
}

// verifyPassword проверяет пароль с солью
func VerifyPassword(password, salt, storedHash string) bool {
	hasher := sha512.New()

	// Хешируем пароль + соль
	combined := password + salt
	hasher.Write([]byte(combined))
	computedHashBytes := hasher.Sum(nil)
	computedHash := hex.EncodeToString(computedHashBytes)

	return computedHash == storedHash
}
