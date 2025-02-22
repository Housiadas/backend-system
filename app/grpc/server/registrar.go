package server

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	userV1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
)

func (s *Server) Registrar() *grpc.Server {
	// -------------------------------------------------------------------------
	// Health Server
	// -------------------------------------------------------------------------
	healthServer := health.NewServer()
	go func() {
		for {
			status := healthpb.HealthCheckResponse_SERVING
			// Check if user Service is valid
			if time.Now().Second()%2 == 0 {
				status = healthpb.HealthCheckResponse_NOT_SERVING
			}

			healthServer.SetServingStatus(userV1.UserService_ServiceDesc.ServiceName, status)
			healthServer.SetServingStatus("", status)

			time.Sleep(1 * time.Second)
		}
	}()

	// -------------------------------------------------------------------------
	// Register gRPC services
	// -------------------------------------------------------------------------
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.grpcInterceptor),
	)
	userV1.RegisterUserServiceServer(grpcServer, s)
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	reflection.Register(grpcServer)

	return grpcServer
}
