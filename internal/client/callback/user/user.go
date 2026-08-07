package user

import (
	"context"

	pb "github.com/artni96/GophKeeper/api/proto/users"
	clientapp "github.com/artni96/GophKeeper/internal/client/app"
	"github.com/artni96/GophKeeper/internal/client/model/user"
)

func LoginCallback(ctx context.Context, app *clientapp.App, loginEntity user.LoginRequest) error {
	req := &pb.LoginRequest{}
	req.SetUsername(loginEntity.Login)
	req.SetPassword(loginEntity.Password)

	token, err := app.UserClient.Login(ctx, req)
	if err != nil {
		return err
	}
	app.State.Token = token.GetToken()
	return nil
}

func RegisterCallback(ctx context.Context, app *clientapp.App, loginEntity user.RegisterRequest) error {
	req := &pb.UserCreateRequest{}
	req.SetUsername(loginEntity.Login)
	req.SetPassword(loginEntity.Password)

	_, err := app.UserClient.CreateUser(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
