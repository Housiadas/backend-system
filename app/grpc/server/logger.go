package server

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GrpcInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	if err != nil {
		s.Log.Error(ctx, "Error during grpc request", err)
	}

	defer func() {
		s.Log.Info(ctx, "grpc request completed",
			"method", info.FullMethod,
			"status_code", int(statusCode),
			"status_text", statusCode.String(),
			"user_agent", s.extractMetadata(ctx).UserAgent,
			"execution_time", duration,
		)
	}()

	return result, err
}
