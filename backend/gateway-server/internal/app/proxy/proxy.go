package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
)

func New(cfg *config.Config) (*http.ServeMux, error) {
	// Article
	target, err := url.Parse(sockerFromParams(
		cfg.Server.ArticleService.Protocol, cfg.Server.ArticleService.Host, cfg.Server.Port,
	))
	if err != nil {
		return nil, err
	}
	articleProxy := httputil.NewSingleHostReverseProxy(target)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/articles/", articleProxy)

	return mux, nil
}

func sockerFromParams(
	protocol, host string,
	port int,
) string {
	return fmt.Sprintf(
		"%s://%s:%d",
		protocol,
		host,
		port,
	)
}
