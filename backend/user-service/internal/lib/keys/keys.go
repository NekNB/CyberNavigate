package keys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadKeyPair загружает пару ключей (приватный и публичный) из файлов
func LoadKeyPair(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Загружаем приватный ключ
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Загружаем публичный ключ
	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load public key: %w", err)
	}

	return privateKey, publicKey, nil
}

// loadPrivateKey загружает приватный ключ из PEM файла
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// Читаем файл
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Декодируем PEM блок
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Парсим приватный ключ (PKCS#1 формат)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Пробуем PKCS#8 формат
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not RSA private key")
		}
	}

	return privateKey, nil
}

// loadPublicKey загружает публичный ключ из PEM файла
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	// Читаем файл
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	// Декодируем PEM блок
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Парсим публичный ключ (PKIX формат)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Приводим к RSA публичному ключу
	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA public key")
	}

	return publicKey, nil
}
