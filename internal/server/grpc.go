package server

import (
	"fmt"
	"net"

	"github.com/artni96/GophKeeper/internal/config"
	"github.com/artni96/GophKeeper/internal/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GRPCServer represents gRPC server structure for its initialization.
type GRPCServer struct {
	cfg    *config.Config
	logger *zap.Logger
	server *grpc.Server
	creds  credentials.TransportCredentials
}

// Init initializes a new grpc server instance.
func (s *GRPCServer) Init() error {
	if s.cfg.EnableTCP {
		err := s.PrepareCredentials()
		if err != nil {
			return err
		}

		s.server = grpc.NewServer(
			grpc.Creds(s.creds),
			grpc.ChainUnaryInterceptor(
				interceptors.RequestLoggerInterceptor(s.logger),
				//interceptors.AuthInterceptor(s.cfg, s.logger),
				interceptors.PanicInterceptor(s.logger),
			),
		)
		return nil

	}
	s.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestLoggerInterceptor(s.logger),
			//interceptors.AuthInterceptor(s.cfg, s.logger),
			interceptors.PanicInterceptor(s.logger),
		),
	)

	return nil
}

// PrepareCredentials sets tcp credentials for the gRPC server.
func (s *GRPCServer) PrepareCredentials() error {
	creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFile, s.cfg.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to create credentials for gRPC server: %w", err)
	}
	s.creds = creds
	return nil
}

// Launch starts the gRPC server.
func (s *GRPCServer) Launch() error {
	listen, err := net.Listen("tcp", s.cfg.ServerAddr)
	if err != nil {
		return fmt.Errorf("failed to announce on the local network address.: %w", err)
	}

	if err = s.server.Serve(listen); err != nil {
		return fmt.Errorf("failed to launch gRPC server: %w", err)
	}
	return nil
}

// Stop stops the gRPC server gracefully.
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// NewGRPCServer returns a new gRPC server instance.
func NewGRPCServer(cfg *config.Config, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		cfg:    cfg,
		logger: logger,
		//urlService:  urlService,
		//userService: userService,
	}
}
