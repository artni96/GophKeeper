package grpc

import (
	"fmt"
	"net"

	cardspb "github.com/artni96/GophKeeper/api/proto/cards"
	healthpb "github.com/artni96/GophKeeper/api/proto/health"
	loginspb "github.com/artni96/GophKeeper/api/proto/logins"
	textspb "github.com/artni96/GophKeeper/api/proto/texts"
	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/server/config"
	cardhandler "github.com/artni96/GophKeeper/internal/server/handler/card"
	healthhandler "github.com/artni96/GophKeeper/internal/server/handler/heath"
	loginhandler "github.com/artni96/GophKeeper/internal/server/handler/login"
	texthandler "github.com/artni96/GophKeeper/internal/server/handler/text"
	userhandler "github.com/artni96/GophKeeper/internal/server/handler/user"
	"github.com/artni96/GophKeeper/internal/server/interceptors"
	cardserv "github.com/artni96/GophKeeper/internal/server/service/card"
	loginserv "github.com/artni96/GophKeeper/internal/server/service/login"
	textserv "github.com/artni96/GophKeeper/internal/server/service/text"
	userserv "github.com/artni96/GophKeeper/internal/server/service/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GRPCServer represents gRPC server structure for its initialization.
type GRPCServer struct {
	cfg    *config.Config
	logger *zap.Logger
	server *grpc.Server

	creds       credentials.TransportCredentials
	userService *userserv.Service
	userHandler *userhandler.Handler

	loginService *loginserv.Service
	loginHandler *loginhandler.Handler

	cardService *cardserv.Service
	cardHandler *cardhandler.Handler

	textService *textserv.Service
	textHandler *texthandler.Handler

	healthHandler *healthhandler.Handler

	streams map[uuid.UUID][]chan *userspb.UpdateNotification
}

// Init initializes a new gRPC server instance.
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
				interceptors.AuthInterceptor(s.cfg, s.logger),
				interceptors.PanicInterceptor(s.logger),
			),
		)
		return nil

	}
	s.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestLoggerInterceptor(s.logger),
			interceptors.AuthInterceptor(s.cfg, s.logger),
			interceptors.PanicInterceptor(s.logger),
		),
	)

	return nil
}

// PrepareCredentials sets tcp credentials for the gRPC server.
func (s *GRPCServer) PrepareCredentials() error {
	creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFile, s.cfg.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to create credentials for gRPC grpc: %w", err)
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

	userHandler := userhandler.NewHandler(s.userService, s.logger, s.streams, s.cfg)
	userspb.RegisterUserServiceServer(s.server, userHandler)

	loginHandler := loginhandler.NewHandler(s.loginService, s.logger, s.streams)
	loginspb.RegisterLoginServiceServer(s.server, loginHandler)

	cardHandler := cardhandler.NewHandler(s.cardService, s.logger, s.streams)
	cardspb.RegisterCardServiceServer(s.server, cardHandler)

	textHandler := texthandler.NewHandler(s.textService, s.logger, s.streams)
	textspb.RegisterTextServiceServer(s.server, textHandler)

	healthHandler := healthhandler.NewHandler(s.logger)
	healthpb.RegisterHealthServiceServer(s.server, healthHandler)

	if err = s.server.Serve(listen); err != nil {
		return fmt.Errorf("failed to launch gRPC grpc: %w", err)
	}
	return nil
}

// Stop stops the gRPC grpc gracefully.
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// NewGRPCServer returns a new gRPC grpc instance.
func NewGRPCServer(
	cfg *config.Config,
	logger *zap.Logger,
	userService *userserv.Service,
	loginService *loginserv.Service,
	cardService *cardserv.Service,
	textService *textserv.Service,
	streams map[uuid.UUID][]chan *userspb.UpdateNotification,
) *GRPCServer {
	return &GRPCServer{
		cfg:          cfg,
		logger:       logger,
		userService:  userService,
		loginService: loginService,
		cardService:  cardService,
		textService:  textService,
		streams:      streams,
	}
}
