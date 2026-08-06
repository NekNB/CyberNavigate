package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/NekNB/CyberNavigate/init/internal/assets"
	"github.com/NekNB/CyberNavigate/init/internal/http"
	"github.com/NekNB/CyberNavigate/init/internal/models"
	"github.com/sirupsen/logrus"

	"github.com/yuin/goldmark"
)

func ProcessArticles(apiClient *http.APIClient, log *logrus.Logger) error {

	// 1. Читаем файл articles.json
	jsonPath := path.Join("articles", "articles.json")
	jsonBytes, err := assets.ArticlesFS.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения articles.json: %w", err)
	}

	// 2. Парсим JSON в слайс структур
	var metaList []models.ArticleMeta
	if err := json.Unmarshal(jsonBytes, &metaList); err != nil {
		return fmt.Errorf("ошибка парсинга articles.json: %w", err)
	}

	// 3. Проходимся по полученному списку
	for _, meta := range metaList {
		// Проверяем, чтобы имя файла из JSON имело расширение .md
		// (на случай, если в JSON забыли его указать)
		fileName := meta.Filename
		if !strings.HasSuffix(fileName, ".md") {
			fileName += ".md"
		}

		// Формируем путь до файла внутри embed.FS
		filePath := path.Join("articles", fileName)

		// 4. Читаем содержимое файла статьи
		contentBytes, err := assets.ArticlesFS.ReadFile(filePath)
		if err != nil {
			// Если файла нет, выводим ошибку, но можно продолжить (continue),
			// чтобы остальные статьи обработались
			return fmt.Errorf("ошибка чтения файла статьи %s: %w", filePath, err)
		}

		// 5. Конвертируем Markdown в HTML
		var htmlBuf bytes.Buffer
		if err := goldmark.Convert(contentBytes, &htmlBuf); err != nil {
			return fmt.Errorf("ошибка конвертации MD в HTML для %s: %w", filePath, err)
		}
		htmlContent := htmlBuf.String()

		resp, err := apiClient.CreateArticle(meta.ArticleTitle)
		if err != nil {
			return fmt.Errorf("ошибка отправки статьи '%s': %w", meta.ArticleTitle, err)
		}
		if resp.Code == 201 {
			if resp.Body == nil {
				log.Error("Body не найден")
				continue
			}
			articleId := resp.Body.Id
			resp, err = apiClient.UploadAcrticleText(articleId, htmlContent)
			if err != nil {
				return fmt.Errorf("ошибка отправки статьи '%s': %w", meta.ArticleTitle, err)
			}
			status := "published"
			resp, err = apiClient.UploadAcrticleData(articleId, &status, nil)
			if err != nil {
				return fmt.Errorf("ошибка отправки статьи '%s': %w", meta.ArticleTitle, err)
			}
		}

	}

	return nil
}
