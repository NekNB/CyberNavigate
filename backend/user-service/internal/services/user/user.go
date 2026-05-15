package article

import (
	"errors"

	"github.com/NekNB/CyberNavigate/backend/user-service/internal/domain/models"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/lib/hash"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/services/session"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/storage"
	"github.com/sirupsen/logrus"
)

var (
	AuthenticationError = errors.New("Authorization Error: Incorrect Login or Password")
)

var _ http.UserServiceInterface = (*UserService)(nil)

type UserDataProvider interface {
	Users() ([]*models.UserDTO, error)
	UserByUsername(username string) (user *models.UserDTO, err error)
	UserByUserId(userId string) (user *models.UserDTO, err error)
	PasswordSaltByUsername(username string) (passwordHash string, salt string, err error)
	NewUser(username, passwordHash, salt string) error
}

type UserService struct {
	log              *logrus.Logger
	userDataProvider UserDataProvider
	sessionService   *session.SessionService
}

func New(
	log *logrus.Logger,
	userDataProvider UserDataProvider,
	sessionService *session.SessionService,
) *UserService {
	return &UserService{
		log:              log,
		userDataProvider: userDataProvider,
		sessionService:   sessionService,
	}
}

// Возвращает список Всех пользователей
func (u *UserService) AllUsers() ([]*models.UserDTO, error) {
	return u.userDataProvider.Users()
}

func (u *UserService) Login(username, password string) (*models.TokensDTO, error) {
	storedHash, salt, err := u.userDataProvider.PasswordSaltByUsername(username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, AuthenticationError
		}
		u.log.Error(err)
		return nil, err
	}

	// Проверяем валидность пароля
	if !hash.VerifyPassword(password, salt, storedHash) {
		return nil, AuthenticationError
	}

	user, err := u.userDataProvider.UserByUsername(username)
	if err != nil {
		return nil, err
	}

	return u.sessionService.CreateSession(user.UserId, user.IsAdmin)
}

func (u *UserService) Logout(sessionId string) error {
	return u.sessionService.RevokeSession(sessionId)
}

func (u *UserService) Register(username, password string) error {
	hashPassword, salt, err := hash.HashPasswordWithSalt(password)
	if err != nil {
		u.log.Error(err)
		return err
	}

	return u.userDataProvider.NewUser(username, hashPassword, salt)
}

func (u *UserService) UserByUserId(userId string) (*models.UserDTO, error) {
	return u.userDataProvider.UserByUserId(userId)
}
