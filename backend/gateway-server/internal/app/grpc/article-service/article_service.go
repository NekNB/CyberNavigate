package articleService

import (
	"context"
	"fmt"
	"log"
	"time"

	articlev1 "github.com/NekNB/CyberNavigate/protos/gen/go/neknb.article.v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func CreateConnection(ctx context.Context, serviceSocket string, mux *runtime.ServeMux) error {

	dialCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Новый API
	conn, err := grpc.NewClient(
		serviceSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(loggingUnaryClientInterceptor()),
	)
	if err != nil {
		log.Printf("❌ failed to create gRPC client for %s: %v", serviceSocket, err)
		return err
	}

	// 🔥 ВАЖНО: принудительно проверяем соединение
	if !conn.WaitForStateChange(dialCtx, conn.GetState()) {
		log.Printf("❌ user service not reachable (%s): %v", serviceSocket, err)
		_ = conn.Close()
		return fmt.Errorf("Не удалось установить соединение")
	}

	log.Printf("✅ connected to user service at %s", serviceSocket)

	return articlev1.RegisterArticlesHandler(ctx, mux, conn)
}

func loggingUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)

			log.Printf(
				"❌ gRPC call failed | method=%s | duration=%s | code=%s | err=%v",
				method,
				duration,
				st.Code(),
				st.Message(),
			)
			return err
		}

		log.Printf(
			"✅ gRPC call success | method=%s | duration=%s",
			method,
			duration,
		)

		return nil
	}
}
