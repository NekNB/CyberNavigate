package article

import (
	articlev1 "github.com/NekNB/CyberNavigate/protos/gen/go/neknb.article.v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateGetArticle(req *articlev1.GetArticleRequest) error {
	if req.GetArticleId() == "" {
		return status.Error(codes.InvalidArgument, "article_id is required")
	}
	return nil
}

func validateCreate(req *articlev1.CreateArticleRequest) error {
	if req.GetArticleName() == "" {
		return status.Error(codes.InvalidArgument, "article name is required")
	}

	return nil
}
