package menu

import (
	"context"
	"errors"
	"fmt"
	"os"

	clientapp "github.com/artni96/GophKeeper/internal/client/app"
	"github.com/artni96/GophKeeper/internal/client/callback"
	"github.com/artni96/GophKeeper/internal/client/model"
)

var ErrFailedToImplement = errors.New("failed to implement")

func ServerConnMenu() {

}

func AuthMenu(ctx context.Context, app *clientapp.App) {

	for {
		ok := app.UserState.HasAttempts()
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

		case "2":
			RegisterMenu(ctx, app)
			AuthMenu(ctx, app)
		case "0":
			fmt.Println("Closing client")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Try again.")
			app.UserState.AddAttempt()
			continue
		}
	}
}

func LoginMenu(ctx context.Context, app *clientapp.App) error {
	app.UserState.ResetAttempts()
	loginEntity := model.LoginRequest{}
	for j := range app.UserState.MaxAttempts {
		if j > 0 {
			fmt.Printf("Attempt #%d\n", j+1)
		}

		for i := range app.UserState.MaxAttempts {
			if i > 0 {
				fmt.Printf("Attempt #%d: ", i)
			}
			fmt.Printf("Enter login: ")
			_, err := fmt.Scanln(&loginEntity.Login)
			if err != nil {
				fmt.Println("Invalid login. Try again.")
				app.UserState.AddAttempt()
				i++
				continue
			}
			break
		}

		for i := range app.UserState.MaxAttempts {
			if i > 0 {
				fmt.Printf("Attempt #%d: ", i)
			}
			fmt.Printf("Enter password: ")
			_, err := fmt.Scanln(&loginEntity.Password)
			if err != nil {
				fmt.Println("Invalid password. Try again.")
				app.UserState.AddAttempt()
				i++
				continue
			}
			break
		}

		app.UserState.ResetAttempts()

		err := callback.LoginCallback(ctx, app, loginEntity)
		if err != nil {
			if j < 2 {
				fmt.Println("Failed to login. Try again.")
			} else {
				fmt.Println("Failed to login. Back to menu.")
				fmt.Println()
			}
			j++
			app.UserState.AddAttempt()
			ok := app.UserState.HasAttempts()
			if !ok {
				app.UserState.ResetAttempts()
				return ErrFailedToImplement
			}
			continue
		}
		break
	}
	app.UserState.ResetAttempts()
	return nil
}

func RegisterMenu(ctx context.Context, app *clientapp.App) {
	loginEntity := model.RegisterRequest{}
	fmt.Printf("Enter login: ")
	fmt.Scanln(&loginEntity.Login)
	fmt.Printf("Enter password: ")
	fmt.Scanln(&loginEntity.Password)
	err := callback.RegisterCallback(ctx, app, loginEntity)
	if err != nil {
		return
	}
	fmt.Println("User created successfully!")
	fmt.Println()
}
