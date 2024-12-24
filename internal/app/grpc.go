package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port = flag.Int("port", 50052, "The gRPC server port")

func NewServer(s *grpc.Server) grpc.ServiceRegistrar {
	return grpc.ServiceRegistrar(s)
}

func grpcOptions() fx.Option {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", *port, err)
	}

	return fx.Options(
		fx.Provide(
			grpc.NewServer,
			NewServer,
		),
		fx.Invoke(
			func(s *grpc.Server, lc fx.Lifecycle) {
				reflection.Register(s)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						log.Println("Starting gRPC server on port", *port)

						go func() {
							if err := s.Serve(listener); err != nil {
								log.Fatalf("Failed to start gRPC server: %v", err)
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						log.Println("Stopping gRPC server...")
						s.GracefulStop()
						return nil
					},
				})
			},
		),
	)
}
