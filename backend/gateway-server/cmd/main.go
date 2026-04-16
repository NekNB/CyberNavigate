package main

import (
	"context"
	"log"
	"net/http"

	pb "your_project/gen/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := pb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:50051",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gateway on :8080")
	http.ListenAndServe(":8080", mux)
}
