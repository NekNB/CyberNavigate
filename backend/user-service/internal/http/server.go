package http

import (
	"errors"

	"github.com/NekNB/CyberNavigate/backend/user-service/internal/config"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/domain/models"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/service"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/storage"
	"github.com/NekNB/CyberNavigate/swagger/gen/user"
	"github.com/gofiber/fiber/v3"
	"github.com/oapi-codegen/runtime/types"
	"github.com/sirupsen/logrus"
)

// Здесь реализуем все методы обработки APIServer

// Проверяем, соответствует ли APIServer сгенерированному ServerInterface
var _ user.ServerInterface = (*APIServer)(nil)

type APIServer struct {
	cfg            *config.Config
	log            *logrus.Logger
	userService    UserServiceInterface
	sessionService SessionServiceInterface
}

type UserServiceInterface interface {
	AllUsers() ([]*models.UserDTO, error)
	UserByUserId(userId string) (*models.UserDTO, error)
	Login(username, password string) (*models.TokensDTO, error)
	Register(username, password string) error
}

type SessionServiceInterface interface {
	RefreshAccessToken(refreshToken string) (*models.TokensDTO, error)
	RevokeSession(sessionId string) error
}

func New(
	cfg *config.Config,
	log *logrus.Logger,
	userService UserServiceInterface,
	sessionService SessionServiceInterface) *APIServer {

	return &APIServer{
		cfg:            cfg,
		log:            log,
		userService:    userService,
		sessionService: sessionService,
	}

}

func (a *APIServer) GetAllUsers(c fiber.Ctx) error {
	if c.Locals(config.IsAdminPayload) == nil || !c.Locals(config.IsAdminPayload).(bool) {
		return c.Status(fiber.StatusForbidden).JSON(
			user.ForbiddenResponse{Message: "Access Denied"},
		)
	}

	usersList, err := a.userService.AllUsers()
	if err != nil {
		a.log.Error(err)
		errResponse := fiber.ErrInternalServerError
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}

	var response []*user.UserResponse
	var userId types.UUID
	for _, userData := range usersList {
		userId.Scan(userData.UserId)

		response = append(response, &user.UserResponse{
			Id:        userId,
			Username:  userData.Username,
			IsAdmin:   userData.IsAdmin,
			CreatedAt: &userData.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (a *APIServer) GetCurrentUser(c fiber.Ctx) error {
	if c.Locals(config.UserIdPayload) == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			&user.UnauthorizedResponse{Message: "unauthorized"},
		)
	}

	userData, err := a.userService.UserByUserId(c.Locals(config.UserIdPayload).(string))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				&user.ErrorResponse{Message: err.Error()},
			)
		}
		a.log.Error(err)
		errResponse := fiber.ErrInternalServerError
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}
	var userId types.UUID
	userId.Scan(userData.UserId)

	response := &user.UserResponse{
		Id:        userId,
		Username:  userData.Username,
		IsAdmin:   userData.IsAdmin,
		CreatedAt: &userData.CreatedAt,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (a *APIServer) Login(c fiber.Ctx) error {

	var req *user.UserRequest

	if err := c.Bind().Body(&req); err != nil {
		a.log.Error(err)
		errResponse := fiber.ErrBadRequest
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}
	response, err := a.userService.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.AuthenticationError) {
			return c.Status(fiber.StatusUnauthorized).JSON(
				&user.UnauthorizedResponse{Message: err.Error()},
			)
		}
		a.log.Error(err)
		errResponse := fiber.ErrInternalServerError
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    response.AccessToken,
		MaxAge:   a.cfg.Tokens.Access.Expiration,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    response.RefreshToken,
		MaxAge:   a.cfg.Tokens.Refresh.Expiration,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(&user.MessageResponse{Message: "Login Success"})
}

func (a *APIServer) Logout(c fiber.Ctx) error {

	if err := a.sessionService.RevokeSession(c.Locals(config.SessionIdPayload).(string)); err != nil {
		a.log.Error(err)
		errResponse := fiber.ErrInternalServerError
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func (a *APIServer) RegisterNewUser(c fiber.Ctx) error {
	var req *user.UserRequest

	if err := c.Bind().Body(&req); err != nil {
		a.log.Error(err)
		errResponse := fiber.ErrBadRequest
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}

	if err := a.userService.Register(req.Username, req.Password); err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return c.Status(fiber.StatusBadRequest).JSON(
				&user.ErrorResponse{Message: "Пользователь с таким UserName уже существует"},
			)
		}
		a.log.Error(err)
		errResponse := fiber.ErrInternalServerError
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}
	return c.Status(fiber.StatusCreated).JSON(&user.MessageResponse{Message: "Register Success"})
}

func (a *APIServer) RefreshToken(c fiber.Ctx) error {

	refreshToken := c.Req().Cookies("refreshToken")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(&user.UnprocessableEntityResponse{
				Message: "refreshToken Requiered In Cookie",
			})
	}

	response, err := a.sessionService.RefreshAccessToken(refreshToken)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				&user.NotFoundResponse{Message: err.Error()},
			)
		} else if errors.Is(err, storage.ErrRefreshNotValid) {
			return c.Status(fiber.StatusNotFound).JSON(
				&user.UnauthorizedResponse{Message: err.Error()},
			)
		}
		a.log.Error(err)
		errResponse := fiber.ErrInternalServerError
		return c.Status(errResponse.Code).JSON(
			&user.ErrorResponse{Message: errResponse.Message},
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    response.AccessToken,
		MaxAge:   a.cfg.Tokens.Access.Expiration,
		HTTPOnly: true,
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    response.RefreshToken,
		MaxAge:   a.cfg.Tokens.Refresh.Expiration,
		HTTPOnly: true,
		Path:     "/",
	})
	return c.Status(fiber.StatusCreated).JSON(&user.MessageResponse{Message: "Success"})
}
