package menu

import (
	"context"
	"errors"
	"fmt"
	"os"

	clientapp "github.com/artni96/GophKeeper/internal/client/app"
	usercb "github.com/artni96/GophKeeper/internal/client/callback/user"
	"github.com/artni96/GophKeeper/internal/client/model/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrFailedToImplement = errors.New("failed to implement")

func ServerConnMenu() {

}

func AuthMenu(ctx context.Context, app *clientapp.App) {

	for {
		ok := app.State.HasAttempts()
		if !ok {
			fmt.Println("no attempts left")
			os.Exit(0)
		}

		fmt.Println("Choose a command")
		fmt.Println("1. Login")
		fmt.Println("2. Registration")
		fmt.Println("0. Exit")
		fmt.Println()
		fmt.Printf("The command: ")
		var choice string
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println(err)
		}
		switch choice {
		case "1":
			err = LoginMenu(ctx, app)
			if errors.Is(err, ErrFailedToImplement) {

				continue
			}
			err = app.UploadData(ctx)
			if err != nil {
				os.Exit(0)
			}
			app.EG.Go(func() error {
				app.SeekUpdates(ctx)
				return nil
			})
			fmt.Println("test")

		case "2":
			err = RegisterMenu(ctx, app)
			if errors.Is(err, ErrFailedToImplement) {

				continue
			}
			AuthMenu(ctx, app)
		case "0":
			fmt.Println("Closing client")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Try again.")
			app.State.AddAttempt()
			continue
		}
	}
}

func LoginMenu(ctx context.Context, app *clientapp.App) error {
	app.State.ResetAttempts()
	loginEntity := user.LoginRequest{}
	for j := range app.State.MaxAttempts {
		if j > 0 {
			fmt.Printf("Attempt #%d\n", j+1)
		}

		for i := range app.State.MaxAttempts {
			if i > 0 {
				fmt.Printf("Attempt #%d: ", i)
			}
			fmt.Printf("Enter login: ")
			_, err := fmt.Scanln(&loginEntity.Login)
			if err != nil {
				fmt.Println("Invalid login. Try again.")
				app.State.AddAttempt()
				continue
			}
			break
		}

		for i := range app.State.MaxAttempts {
			if i > 0 {
				fmt.Printf("Attempt #%d: ", i)
			}
			fmt.Printf("Enter password: ")
			_, err := fmt.Scanln(&loginEntity.Password)
			if err != nil {
				fmt.Println("Invalid password. Try again.")
				app.State.AddAttempt()
				continue
			}
			break
		}

		//app.UserState.ResetAttempts()

		err := usercb.LoginCallback(ctx, app, loginEntity)
		if err != nil {
			if j < 2 {
				fmt.Println("Failed to login. Try again.")
			} else {
				fmt.Println("Failed to login. Back to menu.")
				fmt.Println()
			}
			//j++
			app.State.AddAttempt()
			ok := app.State.HasAttempts()
			if !ok {
				app.State.ResetAttempts()
				return ErrFailedToImplement
			}
			continue
		}
		break
	}
	app.State.ResetAttempts()
	return nil
}

func RegisterMenu(ctx context.Context, app *clientapp.App) error {
	loginEntity := user.RegisterRequest{}
	for j := range app.State.MaxAttempts {
		if j > 0 {
			fmt.Println()
			fmt.Printf("Attempt #%d\n", j+1)
		}

		for i := range app.State.MaxAttempts {
			if i > 0 {
				fmt.Printf("Attempt #%d: ", i)
			}
			fmt.Printf("Enter login: ")
			_, err := fmt.Scanln(&loginEntity.Login)
			if err != nil {
				fmt.Println("Invalid login. Try again.")
				app.State.AddAttempt()
				continue
			}
			break
		}

		for i := range app.State.MaxAttempts {
			if i > 0 {
				fmt.Printf("Attempt #%d: ", i)
			}
			fmt.Printf("Enter password: ")
			_, err := fmt.Scanln(&loginEntity.Password)
			if err != nil {
				fmt.Println("Invalid password. Try again.")
				app.State.AddAttempt()
				continue
			}
			break
		}

		err := usercb.RegisterCallback(ctx, app, loginEntity)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				code := st.Code()
				message := st.Message()
				switch {
				case j == 2:
					fmt.Println("Failed to login. Back to menu.")
					fmt.Println()
				case code == codes.AlreadyExists:
					fmt.Printf("Failed to register - %s. Try again.\n", message)
				default:
					fmt.Println("Failed to register. Try again.")
				}
			}

			app.State.AddAttempt()
			ok = app.State.HasAttempts()
			if !ok {
				app.State.ResetAttempts()
				return ErrFailedToImplement
			}
			continue
		}
		break
	}
	fmt.Println("User created successfully!")
	fmt.Println()
	app.State.ResetAttempts()
	return nil

}
