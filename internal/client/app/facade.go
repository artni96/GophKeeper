package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/client/config"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	cardrepo "github.com/artni96/GophKeeper/internal/client/repository/card"
	loginrepo "github.com/artni96/GophKeeper/internal/client/repository/login"
	textrepo "github.com/artni96/GophKeeper/internal/client/repository/text"
	cardserv "github.com/artni96/GophKeeper/internal/client/service/card"
	healthserv "github.com/artni96/GophKeeper/internal/client/service/health"
	loginserv "github.com/artni96/GophKeeper/internal/client/service/login"
	textserv "github.com/artni96/GophKeeper/internal/client/service/text"
	userserv "github.com/artni96/GophKeeper/internal/client/service/user"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type App struct {
	EG        *errgroup.Group
	conn      *grpc.ClientConn
	Cfg       *config.Config
	NumberMap map[uint64]string

	CardService   *cardserv.Service
	LoginService  *loginserv.Service
	TextService   *textserv.Service
	HealthService *healthserv.Service
	UserService   *userserv.Service

	State    *config.State
	ClientID uuid.UUID
	mu       sync.Mutex
	Reader   *bufio.Reader

	NotificationChan chan common.Notification
	IsBeingUpdated   bool
}

func NewApp(eg *errgroup.Group) *App {
	return &App{
		EG:               eg,
		State:            config.NewState(),
		Cfg:              config.NewConfig(),
		ClientID:         uuid.New(),
		NumberMap:        make(map[uint64]string, 100),
		Reader:           bufio.NewReader(os.Stdin),
		NotificationChan: make(chan common.Notification, 100)}
}

func (app *App) InitDependencies() {
	cardRepo := cardrepo.NewRepository(app.NumberMap)
	app.CardService = cardserv.NewService(app.Cfg, app.conn, cardRepo)

	loginRepo := loginrepo.NewRepository(app.NumberMap)
	app.LoginService = loginserv.NewService(app.Cfg, app.conn, loginRepo)

	textRepo := textrepo.NewRepository(app.NumberMap)
	app.TextService = textserv.NewService(app.Cfg, app.conn, textRepo)

	app.UserService = userserv.NewService(app.Cfg, app.State, app.conn, app.NotificationChan, app.IsBeingUpdated)

}

func (app *App) InitServerConn(ctx context.Context) error {
	for i := range app.Cfg.MaxAttempts {
		if ok := app.State.HasAttempts(); !ok {
			fmt.Println("Failed to connect to the server. Closing the client.")
			os.Exit(0)
		}
		if i != 0 {
			fmt.Printf("Attempt #%d\n", i+1)
		}
		fmt.Printf("Enter server address: ")
		_, err := fmt.Scanln(&app.Cfg.ServerAddress)
		if err != nil {
			if err.Error() == "unexpected newline" {
				app.Cfg.ServerAddress = ":3200"
			} else {
				app.State.AddAttempt()
				continue
			}
		}
		fmt.Printf("Trying to connect to the server: %s\n", app.Cfg.ServerAddress)
		conn, err := grpc.NewClient(app.Cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("Failed to connect to the server")
			app.State.AddAttempt()
			continue
		}

		app.conn = conn

		err = app.testServerConn(ctx)
		if err != nil {
			app.State.AddAttempt()
			continue
		}
		app.State.ResetAttempts()
		return nil
	}
	return errors.New("failed to connect to the server")
}

func (app *App) CloseServerConn() {
	if app.conn != nil {
		app.conn.Close()
	}
}

func (app *App) testServerConn(ctx context.Context) error {
	connCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	app.HealthService = healthserv.NewService(app.conn)
	err := app.HealthService.Check(connCtx)
	if err != nil {
		return errors.New("failed to connect to the server")
	}
	fmt.Println("Successfully connected to the server")
	return nil
}

func (app *App) UploadData(ctx context.Context) error {
	reqCtx := app.PrepareMDContext(ctx)
	if err := app.uploadLogins(reqCtx); err != nil {
		return err
	}
	if err := app.uploadCards(reqCtx); err != nil {
		return err
	}
	if err := app.uploadTexts(reqCtx); err != nil {
		return err
	}
	return nil
}

func (app *App) uploadCards(ctx context.Context) error {
	mdCtx := app.PrepareMDContext(ctx)
	err := app.CardService.AddBatch(mdCtx)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) uploadLogins(ctx context.Context) error {
	mdCtx := app.PrepareMDContext(ctx)
	err := app.LoginService.AddBatch(mdCtx)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) uploadTexts(ctx context.Context) error {
	mdCtx := app.PrepareMDContext(ctx)
	err := app.TextService.AddBatch(mdCtx)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) PrepareMDContext(ctx context.Context) context.Context {
	md := metadata.Pairs("authorization", app.Cfg.Token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}

func (app *App) ReadLine() (string, error) {
	line, err := app.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func (app *App) UpdateStorage(ctx context.Context, n common.Notification) error {
	mdCtx := app.PrepareMDContext(ctx)
	entityNumber := n.EntityNumber
	if n.EntityType == userspb.EntityType_Card {
		if n.ActionType == userspb.ActionType_Create {
			err := app.CardService.Add(mdCtx, entityNumber)
			if err != nil {
				return err
			}
		} else if n.ActionType == userspb.ActionType_Update {
			err := app.CardService.Update(mdCtx, entityNumber)
			if err != nil {
				return err
			}
		} else if n.ActionType == userspb.ActionType_Delete {
			err := app.CardService.Delete(entityNumber)
			if err != nil {
				return err
			}
		}
	} else if n.EntityType == userspb.EntityType_Login {
		if n.ActionType == userspb.ActionType_Create {
			err := app.LoginService.Add(mdCtx, entityNumber)
			if err != nil {
				return err
			}
		} else if n.ActionType == userspb.ActionType_Update {
			err := app.LoginService.Update(mdCtx, entityNumber)
			if err != nil {
				return err
			}
		} else if n.ActionType == userspb.ActionType_Delete {
			err := app.LoginService.Delete(entityNumber)
			if err != nil {
				return err
			}
		}
	}
	app.IsBeingUpdated = false
	return nil
}
