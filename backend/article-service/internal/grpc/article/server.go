package article

import (
	"context"
	"errors"

	"github.com/NekNB/ArticleService/sso/internal/storage"
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
		return nil, err
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

func validateCreate(req *articlev1.CreateArticleRequest) error {
	if req.GetArticleName() == "" {
		return status.Error(codes.InvalidArgument, "article name is required")
	}

	return nil
}

func (s *serverAPI) GetArticle(
	ctx context.Context,
	req *articlev1.GetArticleRequest,
) (*articlev1.ArticleStatusResponse, error) {
	if err := validateGetArticle(req); err != nil {
		return nil, err
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

func validateGetArticle(req *articlev1.GetArticleRequest) error {
	if req.GetArticleId() == "" {
		return status.Error(codes.InvalidArgument, "article_id is required")
	}
	return nil
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

func validateIsAdmin(req *articlev1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
