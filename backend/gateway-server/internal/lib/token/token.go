package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID    string `json:"userId"`
	IsAdmin   bool   `json:"isAdmin"`
	SessionID string `json:"sessionId"`
	jwt.RegisteredClaims
}

// ParseAccessToken парсит и валидирует access токен
func ParseAccessToken(tokenString string, publicKey *rsa.PublicKey) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	return claims, nil
}

// GetPublicKeyFromString загружает публичный ключ из PEM строки
func GetPublicKeyFromString(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Пробуем распарсить как PKIX
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// Пробуем как PKCS1
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			// Пробуем как сертификат
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse public key: %w", err)
			}
			pub = cert.PublicKey
		}
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPub, nil
}

// GetPublicKeyFromFile загружает публичный ключ из файла
func GetPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	return GetPublicKeyFromString(string(keyData))
}
