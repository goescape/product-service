package routes

import (
	"log"
	"net"
	productHandler "product-svc/handlers/product"
	"product-svc/middlewares"
	"product-svc/proto/product"
	"sync"

	"google.golang.org/grpc"
)

type Routes struct {
	GrpcServer *grpc.Server
	Listener   net.Listener

	ProductHandler *productHandler.Handler
}

func (r *Routes) RunGRPC(port string, wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	log.Default().Printf("gRPC server is starting on port %s", port)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares.UnaryLoggingInterceptor),
	)

	r.GrpcServer = grpcServer
	r.Listener = listener

	// Register your gRPC services here
	product.RegisterProductServiceServer(grpcServer, r.ProductHandler)

	err = r.GrpcServer.Serve(r.Listener)
	if err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
