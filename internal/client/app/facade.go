package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	cardspb "github.com/artni96/GophKeeper/api/proto/cards"
	healthpb "github.com/artni96/GophKeeper/api/proto/health"
	loginspb "github.com/artni96/GophKeeper/api/proto/logins"
	textspb "github.com/artni96/GophKeeper/api/proto/texts"
	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/client/user"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	eg           *errgroup.Group
	conn         *grpc.ClientConn
	UserClient   userspb.UserServiceClient
	LoginClient  loginspb.LoginServiceClient
	CardClient   cardspb.CardServiceClient
	TextClient   textspb.TextServiceClient
	HealthClient healthpb.HealthServiceClient
	ServerAddr   string
	UserState    *user.UserState
}

func NewApp(eg *errgroup.Group) *App {
	return &App{
		eg:        eg,
		UserState: user.NewUserState(),
	}
}

func (app *App) InitServerConn(ctx context.Context) error {
	for i := range app.UserState.MaxAttempts {
		if ok := app.UserState.HasAttempts(); !ok {
			fmt.Println("Failed to connect to the server. Closing the client.")
			os.Exit(0)
		}
		if i != 0 {
			fmt.Printf("Attempt #%d\n", i+1)
		}
		fmt.Printf("Enter server address: ")
		_, err := fmt.Scanln(&app.ServerAddr)
		if err != nil {
			if err.Error() == "unexpected newline" {
				app.ServerAddr = ":3200"
			} else {
				app.UserState.AddAttempt()
				continue
			}
		}
		fmt.Printf("Trying to connect to the server: %s\n", app.ServerAddr)
		conn, err := grpc.NewClient(app.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("Failed to connect to the server")
			app.UserState.AddAttempt()
			continue
		}

		app.conn = conn

		err = app.testServerConn(ctx)
		if err != nil {
			app.UserState.AddAttempt()
			continue
		}
		app.UserState.ResetAttempts()
		return nil
	}
	return errors.New("failed to connect to the server")
}

func (app *App) CloseServerConn() {
	if app.conn != nil {
		app.conn.Close()
	}
}

func (app *App) InitClients() {
	app.UserClient = userspb.NewUserServiceClient(app.conn)
	app.LoginClient = loginspb.NewLoginServiceClient(app.conn)
	app.CardClient = cardspb.NewCardServiceClient(app.conn)
	app.TextClient = textspb.NewTextServiceClient(app.conn)
}

func (app *App) testServerConn(ctx context.Context) error {
	connCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	app.HealthClient = healthpb.NewHealthServiceClient(app.conn)
	_, err := app.HealthClient.CheckHealth(connCtx, &healthpb.HealthCheckRequest{})
	if err != nil {
		return errors.New("failed to connect to the server")
	}
	fmt.Println("Successfully connected to the server")
	return nil
}
