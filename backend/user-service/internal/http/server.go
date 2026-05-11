package http

import (
	"github.com/NekNB/CyberNavigate/swagger/gen/user"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

// Здесь реализуем все методы обработки APIServer

// Проверяем, соответствует ли APIServer сгенерированному ServerInterface
var _ user.ServerInterface = (*APIServer)(nil)

type APIServer struct {
	log            *logrus.Logger
	articleService ArticleServiceInterface
}

type ArticleServiceInterface interface {
}

func New(log *logrus.Logger, articleService ArticleServiceInterface) *APIServer {

	return &APIServer{log: log, articleService: articleService}
}

func (a *APIServer) GetAllUsers(c fiber.Ctx) error {

}
