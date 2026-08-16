package health

import (
	"context"

	pb "github.com/artni96/GophKeeper/api/proto/health"
	"go.uber.org/zap"
)

// Handler represents the Health handler instance.
type Handler struct {
	pb.UnimplementedHealthServiceServer
	Logger *zap.Logger
}

// NewHandler initializes and returns the Health Handler instance.
func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{Logger: logger}
}

func (h *Handler) CheckHealth(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{}, nil
}
