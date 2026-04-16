package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"

	pb "github.com/NekNB/CyberNavigate/protos/gen/user/neknb.user.v1"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/assets"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mainMux := http.NewServeMux()

	// 1. Настройка раздачи JSON файлов

	// Открываем подсистему папки docs
	specsDir, err := fs.Sub(assets.SpecsFS, "docs/proto")
	if err != nil {
		panic(err)
	}
	specServer := http.FileServer(http.FS(specsDir))

	// Вешаем их на адреса /specs/...
	// Теперь файлы доступны по URL:
	// http://localhost:8080/specs/user.swagger.json
	// http://localhost:8080/specs/article.swagger.json
	mainMux.Handle("/specs/", http.StripPrefix("/specs/", specServer))

	// 2. Настройка раздачи самого интерфейса (как в прошлом ответе)
	uiDir, _ := fs.Sub(assets.SwaggerUI, "swagger")
	uiServer := http.FileServer(http.FS(uiDir))
	mainMux.Handle("/swagger/", http.StripPrefix("/swagger/", uiServer))

	mux := runtime.NewServeMux()

	err = pb.RegisterUserHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:50051",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	mainMux.Handle("/api/", mux)

	log.Println("gateway on :8080")

	// важно: ты сейчас НЕ используешь mux в ListenAndServe
	http.ListenAndServe(":8080", mainMux)
}
