package article

import (
	"context"
	"errors"

	"github.com/NekNB/CyberNavigate/backend/article-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/storage"
	"github.com/NekNB/CyberNavigate/swagger/gen/article"
	"github.com/sirupsen/logrus"
)

var _ http.ArticleServiceInterface = (*ArticleService)(nil)

type ArticleContentProvider interface {
	ArticleTextById(ctx context.Context, id string) (string, error)
	SaveText(ctx context.Context, articleText string) (string, error)
	UpdateText(ctx context.Context, id string, articleText string) error
}

type ArticleMetaProvider interface {
	Articles() (*[]article.ArticleMetaData, error)
	ArticleByUUID(articleUUID string) (*article.ArticleMetaData, error)
	ArticleTextIDByUUID(articleUUID string) (textID string, err error)
	CreateArticle(articleName *string) (*article.ArticleMetaData, error)
	UpdateArticleByUUID(articleUUID string, title, textID, status, videoUrl *string) (*article.ArticleMetaData, error)
}

type ArticleService struct {
	log                    *logrus.Logger
	articleContentProvider ArticleContentProvider
	articleMetaProvider    ArticleMetaProvider
}

func New(
	log *logrus.Logger,
	articleContentProvider ArticleContentProvider,
	articleMetaProvider ArticleMetaProvider,
) *ArticleService {
	return &ArticleService{
		log:                    log,
		articleContentProvider: articleContentProvider,
		articleMetaProvider:    articleMetaProvider,
	}
}

// Возвращает списко MetaData статей
func (a *ArticleService) Articles() (*[]article.ArticleMetaData, error) {
	metadata, err := a.articleMetaProvider.Articles()
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}

func (a *ArticleService) ArticleByUUID(articleId string) (*article.ArticleMetaData, error) {
	metadata, err := a.articleMetaProvider.ArticleByUUID(articleId)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}

func (a *ArticleService) ArticleTextByUUID(ctx context.Context, articleId string) (string, error) {
	mP := a.articleMetaProvider
	cP := a.articleContentProvider

	textId, err := mP.ArticleTextIDByUUID(articleId)
	if err != nil {
		if !errors.Is(err, storage.ErrArticleTextNotCreatedYet) {
			a.log.Error(err)
		}

		return "", err
	}

	text, err := cP.ArticleTextById(ctx, textId)
	if err != nil {
		a.log.Error(err)
		return text, err
	}

	return text, nil
}

func (a *ArticleService) CreateArticle(title *string) (*article.ArticleMetaData, error) {
	metadata, err := a.articleMetaProvider.CreateArticle(title)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, err
}

func (a *ArticleService) SaveArticleTextByUUID(ctx context.Context, articleId, text string) (*article.ArticleMetaData, error) {
	mP := a.articleMetaProvider
	cP := a.articleContentProvider

	textId, err := cP.SaveText(ctx, text)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}

	metadata, err := mP.UpdateArticleByUUID(
		articleId,
		nil,
		&textId,
		nil,
		nil,
	)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}

func (a *ArticleService) UpdateArticleByUUID(articleId string, title, status *string) (*article.ArticleMetaData, error) {
	metadata, err := a.articleMetaProvider.UpdateArticleByUUID(
		articleId,
		title,
		nil,
		status,
		nil,
	)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, err
}

func (a *ArticleService) UpdateArticleTextByUUID(ctx context.Context, articleId, text string) (*article.ArticleMetaData, error) {
	mP := a.articleMetaProvider
	cP := a.articleContentProvider

	textId, err := mP.ArticleTextIDByUUID(articleId)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	a.log.Debug(textId)
	if err := cP.UpdateText(ctx, textId, text); err != nil {
		a.log.Error(err)
		return nil, err
	}

	metadata, err := mP.ArticleByUUID(
		articleId,
	)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}
