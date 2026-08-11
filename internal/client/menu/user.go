package menu

import (
	"context"
	"fmt"
	"strings"

	"github.com/artni96/GophKeeper/internal/client/model/user"
)

func (m *Menu) initUserMenu() {
	m.routes[register] = StepInfo{
		Logic: func(ctx context.Context) error {
			userEntity := user.RegistrationRequest{}
			for i := range m.app.Cfg.MaxAttempts {
				fmt.Print("Enter login: ")
				loginVal, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid login")
					continue
				}
				userEntity.Login = strings.TrimSpace(loginVal)
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Print("Enter password: ")
				passwordVal, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid password")
					continue
				}
				userEntity.Password = strings.TrimSpace(passwordVal)
				break
			}

			err := m.app.UserService.Register(ctx, userEntity)
			if err != nil {
				fmt.Println("Failed to register user")

				return err
			}
			fmt.Println("User registered successfully!")
			fmt.Println()

			fmt.Print()
			fmt.Println("1. Register")
			fmt.Println("2. Login")
			fmt.Println("0. Exit")
			fmt.Print("Choose option: ")
			m.needInput = true
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			switch choice {
			case "1":
				return register, nil
			case "2":
				return login, nil
			case "0":
				return exit, nil
			default:
				fmt.Println("Invalid choice. Try again")
				return initial, nil
			}
		},
		NextSteps: []int{login, register, exit},
	}

	m.routes[login] = StepInfo{
		Logic: func(ctx context.Context) error {
			userEntity := user.LoginRequest{}
			for i := range m.app.Cfg.MaxAttempts {
				fmt.Print("Enter login: ")
				loginVal, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid login")
					continue
				}
				userEntity.Login = strings.TrimSpace(loginVal)
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Print("Enter password: ")
				passwordVal, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid password")
					continue
				}
				userEntity.Password = strings.TrimSpace(passwordVal)
				break
			}

			err := m.app.UserService.Login(ctx, userEntity)
			if err != nil {
				fmt.Println("Failed to login.")
				fmt.Println()
				m.isFailed = true
				m.needInput = false
				return err
			}

			m.app.State.IsOnline = true
			m.needInput = false

			err = m.app.UploadData(ctx)
			if err != nil {
				return err
			}
			m.app.IsBeingUpdated = true
			fmt.Println("Successfully logged in!")
			fmt.Println()

			m.app.EG.Go(func() error {
				m.app.UserService.SeekUpdates(ctx)
				return nil
			})

			m.app.EG.Go(func() error {
				for {
					select {
					case n, ok := <-m.app.NotificationChan:
						if ok {
							err = m.app.UpdateStorage(ctx, n)
							if err != nil {
								fmt.Println("DEBUG - Failed to update notification")
							}
						} else {
							return nil
						}
					}
				}
			})
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			if m.app.State.IsOnline {
				return main, nil
			}
			if m.isFailed {
				m.isFailed = false
				return initial, nil
			}
			return main, nil
		},
		NextSteps: []int{main, login, exit},
	}
}
