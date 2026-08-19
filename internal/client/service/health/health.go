package health

import (
	"context"

	"github.com/artni96/GophKeeper/api/proto/health"
	"google.golang.org/grpc"
)

type Service struct {
	client health.HealthServiceClient
}

func NewService(conn *grpc.ClientConn) *Service {
	return &Service{
		client: health.NewHealthServiceClient(conn),
	}
}

func (s *Service) Check(ctx context.Context) error {
	req := health.HealthCheckRequest{}
	_, err := s.client.CheckHealth(ctx, &req)
	return err
}
