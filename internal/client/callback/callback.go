package callback

import (
	"context"
	"fmt"

	pb "github.com/artni96/GophKeeper/api/proto/users"
	clientapp "github.com/artni96/GophKeeper/internal/client/app"
	"github.com/artni96/GophKeeper/internal/client/model"
)

func LoginCallback(ctx context.Context, app *clientapp.App, loginEntity model.LoginRequest) error {
	req := &pb.LoginRequest{}
	req.SetUsername(loginEntity.Login)
	req.SetPassword(loginEntity.Password)

	token, err := app.UserClient.Login(ctx, req)
	if err != nil {
		return err
	}
	app.UserState.Token = token.GetToken()
	return nil
}

func RegisterCallback(ctx context.Context, app *clientapp.App, loginEntity model.RegisterRequest) error {
	req := &pb.UserCreateRequest{}
	req.SetUsername(loginEntity.Login)
	req.SetPassword(loginEntity.Password)

	_, err := app.UserClient.CreateUser(ctx, req)
	if err != nil {
		fmt.Println("Failed to create user")
		if err.Error() == "User already exists" {
			return err
		}
		return err
	}
	return nil
}
