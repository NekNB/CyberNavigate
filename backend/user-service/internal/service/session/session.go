package session

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/NekNB/CyberNavigate/backend/user-service/internal/config"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/domain/models"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/lib/keys"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type SessionService struct {
	cfg             *config.Config
	log             *logrus.Logger
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	sessionProvider SessionProvider
}

type AccessClaims struct {
	UserID    string `json:"userId"`
	IsAdmin   bool   `json:"isAdmin"`
	SessionID string `json:"sessionId"`
	jwt.RegisteredClaims
}

var SessionServerError = &models.CustomError{}

type SessionProvider interface {
	ExtendSession(refreshToken string, expiration int) (session_id string, err error)
	NewSession(userId string, refreshTokenExpiration int) (session *models.SessionDTO, err error)
	UserInfoByRefreshToken(refreshToken string) (user *models.UserDTO, err error)
	RevokeSession(sessionId string) error
}

var _ http.SessionServiceInterface = (*SessionService)(nil)

func New(cfg *config.Config, log *logrus.Logger, sessionProvider SessionProvider) (*SessionService, error) {
	privateKey, publicKey, err := keys.LoadKeyPair(cfg.PrivateKeyPath, cfg.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	return &SessionService{
		cfg:             cfg,
		log:             log,
		privateKey:      privateKey,
		publicKey:       publicKey,
		sessionProvider: sessionProvider,
	}, nil
}

// GenerateAccessToken генерирует access token
func (s *SessionService) GenerateAccessToken(userID string, isAdmin bool, sessionID string) (string, error) {
	claims := AccessClaims{
		UserID:    userID,
		IsAdmin:   isAdmin,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.Tokens.Access.Expiration) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

// GenerateRefreshToken генерирует refresh token
func (s *SessionService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 8) // 8 байт = 16 символов в hex
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

// ParseAndVerifyToken парсит и проверяет токен
func (s *SessionService) ParseAndVerifyAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
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

func (s *SessionService) RefreshAccessToken(refreshToken string) (*models.TokensDTO, error) {
	// Продлевает действие refresh, если он не истек и если он существует
	sessionId, err := s.sessionProvider.ExtendSession(refreshToken, s.cfg.Tokens.Refresh.Expiration)
	if err != nil {
		if !errors.Is(err, storage.ErrRefreshNotValid) {
			s.log.Error(err)
		}
		return nil, err
	}

	// Получаем данные пользователя, чтобы создать новый refresh_token
	user, err := s.sessionProvider.UserInfoByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if user.UserId == "" || user.IsAdmin == false {
		err := &models.CustomError{Msg: fmt.Sprintf("Невалидные данные: %+v", user)}
		s.log.Error(err)
		return nil, err
	}

	accessToken, err := s.GenerateAccessToken(user.UserId, user.IsAdmin, sessionId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return &models.TokensDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *SessionService) CreateSession(userId string, isAdmin bool) (*models.TokensDTO, error) {
	session, err := s.sessionProvider.NewSession(userId, s.cfg.Tokens.Refresh.Expiration)
	if err != nil {
		if !errors.Is(err, storage.ErrUserNotFound) {
			s.log.Error(err)
		}
		return nil, err
	}
	if session.SessionId == "" || session.RefreshToken == "" {
		err := &models.CustomError{Msg: fmt.Sprintf("Невалидные данные: %+v", session)}
		s.log.Error(err)
		return nil, err
	}

	accessToken, err := s.GenerateAccessToken(userId, isAdmin, session.SessionId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return &models.TokensDTO{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (s *SessionService) RevokeSession(sessionId string) error {
	return s.sessionProvider.RevokeSession(sessionId)
}
