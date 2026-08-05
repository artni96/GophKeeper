package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/artni96/GophKeeper/internal/config"
	applog "github.com/artni96/GophKeeper/internal/logger"
	cardrepo "github.com/artni96/GophKeeper/internal/repository/card"
	loginrepo "github.com/artni96/GophKeeper/internal/repository/login"
	"github.com/artni96/GophKeeper/internal/repository/text"
	userrepo "github.com/artni96/GophKeeper/internal/repository/user"
	"github.com/artni96/GophKeeper/internal/server"
	cardserv "github.com/artni96/GophKeeper/internal/service/card"
	loginserv "github.com/artni96/GophKeeper/internal/service/login"
	textserv "github.com/artni96/GophKeeper/internal/service/text"
	userserv "github.com/artni96/GophKeeper/internal/service/user"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// App is the app for the app with all dependencies.
type App struct {
	eg           *errgroup.Group
	Cfg          *config.Config
	DB           *sqlx.DB
	Logger       *zap.Logger
	server       server.GRPCServer
	userService  *userserv.Service
	loginService *loginserv.Service
	cardService  *cardserv.Service
	textService  *textserv.Service
}

// NewApp initializes and returns a new instance of App.
func NewApp(eg *errgroup.Group, cfg *config.Config) *App {
	return &App{
		eg:  eg,
		Cfg: cfg,
	}
}

// initLogger initializes the app logger.
func (a *App) initLogger() error {
	logger, err := applog.InitLogger(a.Cfg.LoggingLvl)
	if err != nil {
		return fmt.Errorf("failed to initialize app logger: %w", err)
	}
	a.Logger = logger
	a.Logger.Info("app logger initialized successfully")
	return nil
}

// InitDBConn initializes a database connection according to the app config.
func (a *App) InitDBConn(ctx context.Context) error {
	if a.Cfg.DBDsn == "" {
		return fmt.Errorf("database dsn is not provided")
	}

	db, err := sqlx.Open("pgx", a.Cfg.DBDsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	localCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err = db.PingContext(localCtx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	a.DB = db
	a.Logger.Info("database connection initialized successfully")
	return nil
}

// applyMigrations updates the atabase up to the latest migration file.
func (a *App) applyMigrations() error {
	driver, err := postgres.WithInstance(a.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize postgres driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	a.Logger.Info("migrations applied successfully")
	return nil
}

// CloseDBConn closes the database connection.
func (a *App) CloseDBConn() error {
	err := a.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

// initDependencies initializes all necessary dependencies - repositories and services.
func (a *App) initDependencies() error {
	userRepository, err := userrepo.NewRepository(a.DB, a.Logger)
	if err != nil {
		a.Logger.Error("failed to initialize user repository", zap.Error(err))
		return fmt.Errorf("failed to initialize user repository: %w", err)
	}
	a.userService = userserv.NewService(a.Cfg, a.Logger, userRepository)

	loginRepository, err := loginrepo.NewRepository(a.DB, a.Logger)
	if err != nil {
		a.Logger.Error("failed to initialize login repository", zap.Error(err))
		return fmt.Errorf("failed to initialize login repository: %w", err)
	}
	a.loginService = loginserv.NewService(a.Cfg, a.Logger, loginRepository)

	cardRepository, err := cardrepo.NewRepository(a.DB, a.Logger)
	if err != nil {
		a.Logger.Error("failed to initialize card repository", zap.Error(err))
		return fmt.Errorf("failed to initialize card repository: %w", err)
	}
	a.cardService = cardserv.NewService(a.Cfg, a.Logger, cardRepository)

	textRepository, err := text.NewRepository(a.DB, a.Logger)
	if err != nil {
		a.Logger.Error("failed to initialize text repository", zap.Error(err))
		return fmt.Errorf("failed to initialize text repository: %w", err)
	}
	a.textService = textserv.NewService(a.Cfg, a.Logger, textRepository)

	return nil
}

// initServer initializes a new gRPC server instance.
func (a *App) initServer() error {
	newServer := server.NewGRPCServer(a.Cfg, a.Logger, a.userService, a.loginService, a.cardService, a.textService)
	err := newServer.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize gRPC server: %w", err)
	}
	a.server = *newServer
	return nil
}

// launchServer starts the gRPC server.
func (a *App) launchServer() error {
	err := a.initServer()
	if err != nil {
		return err
	}

	a.eg.Go(func() error {
		err = a.server.Launch()
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

// Stop stops the running gRPC server and closes the database connection.
func (a *App) Stop(ctx context.Context, ctxCancel context.CancelFunc, isClosedChan chan struct{}) {
	a.eg.Go(func() error {
		a.server.Stop()

		//time.Sleep(36 * time.Second)
		select {
		case <-ctx.Done():
			a.Logger.Info("server stopped during the forced period")
		default:
			a.Logger.Info("server stopped gracefully")
		}

		if err := a.DB.Close(); err != nil {
			a.Logger.Error("failed to close database connection gracefully", zap.Error(err))
			return fmt.Errorf("failed to close database connection gracefully: %w", err)
		}
		select {
		case <-ctx.Done():
			a.Logger.Info("database connection closed during the forceful period")
		default:
			a.Logger.Info("database connection closed gracefully")
		}

		close(isClosedChan)
		ctxCancel()
		return nil
	})
}

// Launch starts the app with all credentials.
func (a *App) Launch(ctx context.Context) error {
	err := a.initLogger()
	if err != nil {
		log.Fatal("failed to initialize logger", zap.Error(err))
		return err
	}

	err = a.InitDBConn(ctx)
	if err != nil {
		a.Logger.Error("failed to initialize database connection", zap.Error(err))
		return err
	}

	err = a.applyMigrations()
	if err != nil {
		a.Logger.Error("failed to apply migrations", zap.Error(err))
		return err
	}

	err = a.initDependencies()
	if err != nil {
		a.Logger.Error("failed to initialize app dependencies", zap.Error(err))
		return err
	}

	err = a.launchServer()
	if err != nil {
		a.Logger.Error("failed to launch server", zap.Error(err))
		return err
	}
	a.Logger.Info("server launched successfully")
	a.Logger.Info(a.configStdout())
	return nil
}

// configStdout builds
func (a *App) configStdout() string {
	var resp string
	resp += fmt.Sprintf("\n\nApp config:\n")
	resp += fmt.Sprintf("	ServerAddr: %s\n", a.Cfg.ServerAddr)
	resp += fmt.Sprintf("	Database name: %s\n", a.Cfg.DBName)
	resp += fmt.Sprintf("	Token expiration period: %f minutes\n", a.Cfg.TokenExp.Minutes())
	resp += fmt.Sprintf("	TCP is enabled: %t\n", a.Cfg.EnableTCP)

	if a.Cfg.CertFile != "" {
		resp += fmt.Sprintf("	TLS certificate is set: %t\n", a.Cfg.CertFile != "")
	}
	if a.Cfg.KeyFile != "" {
		resp += fmt.Sprintf("	TLS key is set: %t\n", a.Cfg.KeyFile != "")
	}
	resp += fmt.Sprintf("	Logging level: %s\n", a.Cfg.LoggingLvl)
	resp += fmt.Sprintf("	Grace period: %d seconds\n", a.Cfg.GPeriod)
	resp += fmt.Sprintf("	Force period: %d seconds\n", a.Cfg.FPeriod)
	return resp
}
