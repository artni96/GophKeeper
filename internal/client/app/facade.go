package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	cardspb "github.com/artni96/GophKeeper/api/proto/cards"
	healthpb "github.com/artni96/GophKeeper/api/proto/health"
	loginspb "github.com/artni96/GophKeeper/api/proto/logins"
	textspb "github.com/artni96/GophKeeper/api/proto/texts"
	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/client/model/card"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	"github.com/artni96/GophKeeper/internal/client/model/login"
	"github.com/artni96/GophKeeper/internal/client/model/text"
	"github.com/artni96/GophKeeper/internal/client/user"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type App struct {
	EG           *errgroup.Group
	conn         *grpc.ClientConn
	UserClient   userspb.UserServiceClient
	LoginClient  loginspb.LoginServiceClient
	CardClient   cardspb.CardServiceClient
	TextClient   textspb.TextServiceClient
	HealthClient healthpb.HealthServiceClient
	ServerAddr   string
	State        *user.State
	ClientID     uuid.UUID
	mu           sync.Mutex
}

func NewApp(eg *errgroup.Group) *App {
	return &App{
		EG:       eg,
		State:    user.NewUserState(),
		ClientID: uuid.New(),
	}
}

func (app *App) InitServerConn(ctx context.Context) error {
	for i := range app.State.MaxAttempts {
		if ok := app.State.HasAttempts(); !ok {
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
				app.State.AddAttempt()
				continue
			}
		}
		fmt.Printf("Trying to connect to the server: %s\n", app.ServerAddr)
		conn, err := grpc.NewClient(app.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (app *App) UploadData(ctx context.Context) error {
	md := metadata.Pairs("authorization", app.State.Token)
	ctxWithMD := metadata.NewOutgoingContext(ctx, md)
	if err := app.uploadLogins(ctxWithMD); err != nil {
		return err
	}
	if err := app.uploadCards(ctxWithMD); err != nil {
		return err
	}
	if err := app.uploadTexts(ctxWithMD); err != nil {
		return err
	}
	app.sortCommonListByNumber()
	return nil
}

func (app *App) sortCommonListByNumber() {
	list := app.State.DataStorage.ShortDataList
	sort.Slice(list, func(i, j int) bool {
		return list[i].Number < list[j].Number
	})
}

func (app *App) uploadLogins(ctx context.Context) error {

	loginListReq := &loginspb.LoginGetListRequest{}

	loginListResp, err := app.LoginClient.GetListLogin(ctx, loginListReq)
	if err != nil {
		return err
	}

	var shortDataList []common.Entity

	pbEntities := loginListResp.GetLogins()

	for _, pbEntity := range pbEntities {
		createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
		if err != nil {
			return err
		}

		var updatedAt time.Time
		updatedAt, err = time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
		if err != nil {
			updatedAt = time.Time{}
		}

		entity := login.Login{
			Login:       pbEntity.GetLogin(),
			Password:    pbEntity.GetPassword(),
			Title:       pbEntity.GetTitle(),
			Number:      pbEntity.GetNumber(),
			URL:         pbEntity.GetUrl(),
			Description: pbEntity.GetDescription(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		app.State.DataStorage.ExtendedDataMap[entity.Number] = entity

		shortDataEntity := common.Entity{
			Title:       pbEntity.GetTitle(),
			Description: pbEntity.GetDescription(),
			Number:      pbEntity.GetNumber(),
			Type:        "login",
		}

		shortDataList = append(shortDataList, shortDataEntity)
	}
	app.State.DataStorage.ShortDataList = append(app.State.DataStorage.ShortDataList, shortDataList...)
	return nil
}

func (app *App) uploadTexts(ctx context.Context) error {
	textListReq := &textspb.TextGetListRequest{}

	textListResp, err := app.TextClient.GetListText(ctx, textListReq)
	if err != nil {
		return err
	}

	var shortDataList []common.Entity

	pbEntities := textListResp.GetTexts()

	for _, pbEntity := range pbEntities {
		createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
		if err != nil {
			return err
		}

		var updatedAt time.Time
		updatedAt, err = time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
		if err != nil {
			updatedAt = time.Time{}
		}

		entity := text.Text{
			Text:        pbEntity.GetText(),
			Title:       pbEntity.GetTitle(),
			Number:      pbEntity.GetNumber(),
			Description: pbEntity.GetDescription(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		app.State.DataStorage.ExtendedDataMap[entity.Number] = entity

		shortDataEntity := common.Entity{
			Title:       pbEntity.GetTitle(),
			Description: pbEntity.GetDescription(),
			Number:      pbEntity.GetNumber(),
			Type:        "text",
		}

		shortDataList = append(shortDataList, shortDataEntity)
	}
	app.State.DataStorage.ShortDataList = append(app.State.DataStorage.ShortDataList, shortDataList...)
	return nil
}

func (app *App) uploadCards(ctx context.Context) error {
	cardListReq := &cardspb.CardGetListRequest{}

	cardListResp, err := app.CardClient.GetListCard(ctx, cardListReq)
	if err != nil {
		return err
	}

	var shortDataList []common.Entity

	pbEntities := cardListResp.GetCards()

	for _, pbEntity := range pbEntities {
		createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
		if err != nil {
			return err
		}

		var updatedAt time.Time
		updatedAt, err = time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
		if err != nil {
			updatedAt = time.Time{}
		}

		entity := card.Card{
			PAN:         pbEntity.GetPan(),
			Holder:      pbEntity.GetHolder(),
			ExpiryDate:  pbEntity.GetExpiryDate(),
			CVV:         pbEntity.GetCvv(),
			PIN:         pbEntity.GetPin(),
			Bank:        pbEntity.GetBank(),
			Brand:       pbEntity.GetBrand(),
			Title:       pbEntity.GetTitle(),
			Number:      pbEntity.GetNumber(),
			Description: pbEntity.GetDescription(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		app.State.DataStorage.ExtendedDataMap[entity.Number] = entity

		shortDataEntity := common.Entity{
			Title:       pbEntity.GetTitle(),
			Description: pbEntity.GetDescription(),
			Number:      pbEntity.GetNumber(),
			Type:        "card",
		}

		shortDataList = append(shortDataList, shortDataEntity)
	}
	app.State.DataStorage.ShortDataList = append(app.State.DataStorage.ShortDataList, shortDataList...)
	return nil
}

func (app *App) SeekUpdates(ctx context.Context) {
	lostConn := false
	md := metadata.Pairs("authorization", app.State.Token)
Reconnection:
	for {
		timeout := 10 * time.Second
		reqCtx := metadata.NewOutgoingContext(ctx, md)
		stream, err := app.UserClient.SeekUpdates(reqCtx, &userspb.SeekUpdateRequest{})
		if err != nil {
			if !lostConn {
				lostConn = true
				app.State.IsOnline = false
				fmt.Println("You are being offline - failed to connect to the server.")
			}
			time.Sleep(timeout)
			continue Reconnection
		}
		if lostConn {
			lostConn = false
			app.State.IsOnline = true
			fmt.Println("You are being online - reconnected to the server.")
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				update, err := stream.Recv()
				if err != nil {
					lostConn = true
					fmt.Println("You are being offline - lost connection to the server.")
					time.Sleep(timeout)
					continue Reconnection
				}
				fmt.Println("update", update)
			}
		}
	}
}
