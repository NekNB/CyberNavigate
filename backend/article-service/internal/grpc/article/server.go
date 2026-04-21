package article

import (
	"context"
	"errors"

	articleModel "github.com/NekNB/CyberNavigate/backend/article-service/internal/domain/models"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/services/article"
	articlev1 "github.com/NekNB/CyberNavigate/protos/gen/go/neknb.article.v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Article interface {
	Create(
		ctx context.Context,
		article_name string,
	) (uuid, name, status string, err error)
	GetArticle(
		ctx context.Context,
		article_id string,
	) (uuid, name, status string, err error)
	GetArticleStream(
		ctx context.Context,
		article_id string,
	) (content byte, err error)
	GetArticles(
		ctx context.Context,
	) (total int64, articles *[]articleModel.Article)
	Update(
		ctx context.Context,
		article_id string,
		title, text, video_link *string,
	) (bool, error)
}

type serverAPI struct {
	articlev1.UnimplementedArticlesServer
	article Article
}

func Register(gRPC *grpc.Server, Article Article) {
	articlev1.RegisterArticlesServer(gRPC, &serverAPI{article: Article})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Create(
	ctx context.Context,
	req *articlev1.CreateArticleRequest,
) (*articlev1.ArticleStatusResponse, error) {

	if err := validateCreate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	articleUuid, articleName, articleStatus, err := s.article.Create(ctx, req.GetArticleName())
	if err != nil {
		if errors.Is(err, article.ErrArticleExists) {
			return nil, status.Error(codes.InvalidArgument, "invalid arguments")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &articlev1.ArticleStatusResponse{
		ArticleUuid:   articleUuid,
		ArticleName:   articleName,
		ArticleStatus: articleStatus,
	}, nil
}

// Возвращает конкретную запись из статьи
func (s *serverAPI) GetArticle(
	ctx context.Context,
	req *articlev1.GetArticleRequest,
) (*articlev1.ArticleStatusResponse, error) {
	if err := validateGetArticle(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	articleUuid, articleName, articleStatus, err := s.article.GetArticle(ctx, req.GetArticleId())
	if err != nil {
		if errors.Is(err, article.ErrArticleNotFound) {
			return nil, status.Error(codes.NotFound, "article not found. Please check article_uuid")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &articlev1.ArticleStatusResponse{
		ArticleUuid:   articleUuid,
		ArticleName:   articleName,
		ArticleStatus: articleStatus,
	}, nil
}

// Передает поток статьи
func (s *serverAPI) GetArticleStream(
	ctx context.Context,
	req *articlev1.GetArticleRequest,
) (*articlev1., error) {
	if err := validateGetArticle(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	articleUuid, articleName, articleStatus, err := s.article.GetArticle(ctx, req.GetArticleId())
	if err != nil {
		if errors.Is(err, article.ErrArticleNotFound) {
			return nil, status.Error(codes.NotFound, "article not found. Please check article_uuid")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &articlev1.ArticleChunk{
		ArticleUuid:   articleUuid,
		ArticleName:   articleName,
		ArticleStatus: articleStatus,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *articlev1.IsAdminRequest,
) (*articlev1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	isAdmin, err := s.Article.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &articlev1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
