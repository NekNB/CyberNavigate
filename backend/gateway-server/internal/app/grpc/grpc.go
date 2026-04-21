package grpc

import (
	"context"
	"fmt"
	"net/http"

	articleService "github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/grpc/article-service"
	userService "github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/grpc/user-service"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

func Register(ctx context.Context, cfg *config.Config, log *logrus.Logger) (mux *runtime.ServeMux) {
	// TODO: добавить повторные попытки подключение к microservices
	mux = runtime.NewServeMux(
		runtime.WithErrorHandler(func(
			ctx context.Context,
			mux *runtime.ServeMux,
			marshaler runtime.Marshaler,
			w http.ResponseWriter,
			r *http.Request,
			err error,
		) {

			st, _ := status.FromError(err)

			log.Printf(
				"❌ gateway error\npath=%s\nmethod=%s\nhttp_error=%s\ngrpc_code=%s\ngrpc_message=%s",
				r.URL.Path,
				r.Method,
				err.Error(),
				st.Code().String(),
				st.Message(),
			)

			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		}),
	)
	var serverConnectedCounter int

	userServiceSocket := fmt.Sprintf("%s:%d", cfg.Server.UserService.Host, cfg.Server.UserService.Port)
	if err := userService.CreateConnection(ctx, userServiceSocket, mux); err != nil {
		log.Errorf("Ошибка подключения UserService: %s на сокете: %s", err.Error(), userServiceSocket)
	} else {
		serverConnectedCounter++
	}

	articleServiceSocket := fmt.Sprintf("%s:%d", cfg.Server.ArticleService.Host, cfg.Server.ArticleService.Port)
	if err := articleService.CreateConnection(ctx, articleServiceSocket, mux); err != nil {
		log.Errorf("Ошибка подключения ArticleService: %s на сокете: %s", err.Error(), articleServiceSocket)
	} else {
		serverConnectedCounter++
	}

	if serverConnectedCounter != 0 {
		log.Infof("Подключено сервисов: %d", serverConnectedCounter)
	} else {
		log.Warn("Нет подключенных сервисов")
	}
	return
}
