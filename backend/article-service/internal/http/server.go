package http

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NekNB/CyberNavigate/backend/article-service/internal/storage"
	"github.com/NekNB/CyberNavigate/swagger/gen/article"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v3"
)

// Здесь реализуем все методы обработки APIServer

// Проверяем, соответствует ли APIServer сгенерированному ServerInterface
var _ article.ServerInterface = (*APIServer)(nil)

type APIServer struct {
	log            *logrus.Logger
	articleService ArticleServiceInterface
}

type ArticleServiceInterface interface {
	CreateArticle(title *string) (*article.ArticleMetaData, error)
	Articles() (*[]article.ArticleMetaData, error)
	ArticleByUUID(articleId string) (*article.ArticleMetaData, error)
	SaveArticleTextByUUID(ctx context.Context, articleId, text string) (*article.ArticleMetaData, error)
	ArticleTextByUUID(ctx context.Context, articleId string) (string, error)
	UpdateArticleTextByUUID(ctx context.Context, articleId, text string) (*article.ArticleMetaData, error)
	UpdateArticleByUUID(articleId, text, title, videoURL, status string) (*article.ArticleMetaData, error)
}

func New(log *logrus.Logger, articleService ArticleServiceInterface) *APIServer {

	return &APIServer{log: log, articleService: articleService}
}

func (a *APIServer) PostArticles(c fiber.Ctx) error {
	aS := a.articleService
	var request article.PostArticlesJSONRequestBody
	if err := c.Bind().Body(&request); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	metadata, err := aS.CreateArticle(request.ArticleTitle)
	if err != nil {
		if errors.Is(err, storage.ErrArticleExists) {
			errMsg := fmt.Sprintf("Article With %s Already Exists", *request.ArticleTitle)
			return c.Status(fiber.StatusBadRequest).
				JSON(article.ArticleAlreadyExists{
					Message: &errMsg,
				})
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(metadata)
}

func (a *APIServer) GetArticles(c fiber.Ctx) error {
	aS := a.articleService

	metadata, err := aS.Articles()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(metadata)
}

func (a *APIServer) GetArticleById(c fiber.Ctx, articleId string) error {
	aS := a.articleService

	metadata, err := aS.ArticleByUUID(articleId)
	if err != nil {
		if errors.Is(err, storage.ErrArticleNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(metadata)
}

func chunkByWords(s string, wordsPerChunk int) []string {
	words := strings.Fields(s)
	var chunks []string

	for i := 0; i < len(words); i += wordsPerChunk {
		end := i + wordsPerChunk
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}

	return chunks
}

func (a *APIServer) GetArticleTextById(c fiber.Ctx, articleId string) error {
	aS := a.articleService

	text, err := aS.ArticleTextByUUID(c.Context(), articleId)
	if err != nil {
		if errors.Is(err, storage.ErrArticleNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		} else if errors.Is(err, storage.ErrArticleTextNotCreatedYet) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("Content-Type", "application/x-ndjson")
	c.Set("Transfer-Encoding", "chunked")

	return c.SendStreamWriter(func(w *bufio.Writer) {
		chunks := chunkByWords(text, 50)

		for _, chunk := range chunks {
			fmt.Fprintln(w, chunk) // важно: \n для NDJSON
			w.Flush()              // отправляем сразу клиенту
			time.Sleep(500 * time.Millisecond)
		}
	})
}

func (a *APIServer) PatchArticleById(c fiber.Ctx, articleId string) error {
	aS := a.articleService

	var request article.PatchArticleByIdJSONRequestBody

	if err := c.Bind().Body(&request); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// TODO: проверить работу
	metadata, err := aS.UpdateArticleByUUID(
		articleId,
		"",
		"",
		*request.ArticleStatus,
		*request.ArticleTitle,
	)
	if err != nil {
		if errors.Is(err, storage.ErrArticleNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(metadata)
}

func (a *APIServer) PostArticleTextById(c fiber.Ctx, articleId string) error {
	aS := a.articleService

	text := string(c.Request().Body())

	metadata, err := aS.SaveArticleTextByUUID(c.Context(), articleId, text)
	if err != nil {
		a.log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(metadata)
}

func (a *APIServer) PutArticleTextById(c fiber.Ctx, articleId string) error {
	aS := a.articleService

	text := string(c.Request().Body())

	metadata, err := aS.UpdateArticleTextByUUID(c.Context(), articleId, text)
	if err != nil {
		if errors.Is(err, storage.ErrArticleNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		} else if errors.Is(err, storage.ErrArticleTextNotCreatedYet) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(metadata)
}
