package server

import (
	"context"
	"log"
	"net"

	"github.com/chirag-bruno/mock-grpc/internal/transport"
	"github.com/chirag-bruno/mock-grpc/pkg/todo"
	"google.golang.org/grpc"
)

type Config struct {
	Mode    transport.Mode
	Address string
}

func Run(cfg Config) error {
	listener, err := transport.NewListener(cfg.Mode, cfg.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Starting gRPC server in %s mode on %s", cfg.Mode, cfg.Address)
	return serve(listener)
}

func serve(listener net.Listener) error {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	todo.RegisterTodoServiceServer(grpcServer, NewTodoServer())

	log.Printf("Server is ready to accept connections")
	return grpcServer.Serve(listener)
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("[gRPC] Method: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("[gRPC] Error: %v", err)
	}
	return resp, err
}
