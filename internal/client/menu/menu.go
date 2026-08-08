package menu

import (
	"context"
	"fmt"
	"os"

	"strings"

	"github.com/artni96/GophKeeper/internal/client/app"
	"github.com/artni96/GophKeeper/internal/client/model/user"
)

const (
	initial = iota
	register
	login
	authenticated
	exit

	cardList
	createCard
	getCard
	updateCard

	listLogin
	createLogin
	getLogin
	updateLogin
)

type StepInfo struct {
	Output    func(ctx *config) error                      // Show the menu
	Handler   func(ctx *config, input string) (int, error) // Process user choice
	NextSteps []int                                        // Available routes from this menu
}

type config struct {
	Token      string
	IsLoggedIn bool
	IsRepeated bool
	NeedChoice bool
	app        *app.App
}

type Menu struct {
	routes map[int]StepInfo
	cfg    *config
}

func InitMenu(app *app.App) *Menu {
	return &Menu{
		routes: make(map[int]StepInfo),
		cfg: &config{
			app: app,
		},
	}
}

func (m *Menu) initSteps(ctx context.Context) {

	m.routes = make(map[int]StepInfo)

	m.initCardList()
	m.initLoginList()

	m.routes[initial] = StepInfo{
		Output: func(cfg *config) error {

			fmt.Println("1. Register")
			fmt.Println("2. Login")
			fmt.Println("0. Exit")
			fmt.Print("Choose option: ")
			cfg.NeedChoice = true
			return nil
		},
		Handler: func(cfg *config, choice string) (int, error) {
			switch choice {
			case "1":
				//m.ctx.IsRepeated = false
				return register, nil
			case "2":
				//m.ctx.IsRepeated = false
				return login, nil
			case "0":
				return exit, nil
			default:
				//m.ctx.IsRepeated = true
				fmt.Println("Invalid choice. Try again")
				return initial, nil
			}
		},
		NextSteps: []int{register, login, exit},
	}

	m.routes[register] = StepInfo{
		Output: func(cfg *config) error {
			userEntity := user.RegistrationRequest{}
			fmt.Print("Enter login: ")
			_, err := fmt.Scanln(&userEntity.Login)
			if err != nil {
				return err
			}
			userEntity.Login = strings.TrimSpace(userEntity.Login)

			fmt.Print("Enter password: ")
			_, err = fmt.Scanln(&userEntity.Password)
			if err != nil {
				return err
			}
			userEntity.Password = strings.TrimSpace(userEntity.Password)

			err = cfg.app.UserService.Register(ctx, userEntity)
			if err != nil {
				fmt.Println("Failed to register user")

				return err
			}
			fmt.Println("User registered successfully!")

			fmt.Print()
			fmt.Println("1. Register")
			fmt.Println("2. Login")
			fmt.Println("0. Exit")
			fmt.Print("Choose option: ")
			cfg.NeedChoice = true
			return nil
		},
		Handler: func(cfg *config, choice string) (int, error) {
			switch choice {
			case "1":
				m.cfg.IsRepeated = false
				return register, nil
			case "2":
				m.cfg.IsRepeated = false
				return login, nil
			case "0":
				return exit, nil
			default:
				m.cfg.IsRepeated = true
				fmt.Println("Invalid choice. Try again")
				return initial, nil
			}
		},
		NextSteps: []int{login, register, exit},
	}

	m.routes[login] = StepInfo{
		Output: func(cfg *config) error {
			fmt.Print("Enter username: ")
			userEntity := user.LoginRequest{}
			_, err := fmt.Scanln(&userEntity.Login)
			if err != nil {
				return err
			}
			userEntity.Login = strings.TrimSpace(userEntity.Login)

			fmt.Print("Enter password: ")
			_, err = fmt.Scanln(&userEntity.Password)
			if err != nil {
				return err
			}
			userEntity.Password = strings.TrimSpace(userEntity.Password)

			err = cfg.app.UserService.Login(ctx, userEntity)
			if err != nil {
				fmt.Println("Failed to login.")
				return err
			}

			cfg.IsLoggedIn = true
			cfg.NeedChoice = false

			err = cfg.app.UploadData(ctx)
			if err != nil {
				return err
			}
			fmt.Println("Successfully logged in!")
			return nil
		},
		Handler: func(cfg *config, choice string) (int, error) {
			if cfg.IsLoggedIn {
				return authenticated, nil
			}
			return login, nil
		},
		NextSteps: []int{authenticated, login, exit},
	}

	m.routes[authenticated] = StepInfo{
		Output: func(cfg *config) error {
			fmt.Printf("Choose option: ")
			fmt.Println("1. Cards")
			fmt.Println("2. Logins")
			fmt.Println("3. Texts")
			fmt.Println("0. Logout & Exit")
			fmt.Print("Choose option: ")
			cfg.NeedChoice = true
			return nil
		},
		Handler: func(cfg *config, choice string) (int, error) {
			switch choice {
			case "1":
				// Call your list function
				//ctx.Reader.ReadString('\n')
				return cardList, nil
			case "2":
				//fmt.Println("Listing Type B items...")
				//ctx.Reader.ReadString('\n')
				return listLogin, nil
			case "3":
				fmt.Println("Listing Type C items...")
				//ctx.Reader.ReadString('\n')
				return authenticated, nil
			case "0":
				cfg.IsLoggedIn = false
				cfg.Token = ""
				fmt.Println("Logged out. Goodbye!")
				return exit, nil
			default:
				fmt.Println("Invalid choice")
				return authenticated, nil
			}
		},
		NextSteps: []int{authenticated, exit},
	}

}

func (m *Menu) Run(ctx context.Context) {
	m.initSteps(ctx)
	currentStep := initial
	failsCounter := 0

	for currentStep != exit {
		route := m.routes[currentStep]

		err := route.Output(m.cfg)
		if err != nil {
			if currentStep == login || currentStep == register {
				failsCounter++
				if failsCounter == 3 {
					if currentStep == login {
						currentStep = initial
						failsCounter = 0
						fmt.Println()
					} else {
						fmt.Println("The client is being closed...")
						os.Exit(0)
					}
				}
			}
			continue
		}

		var choice string
		if m.cfg.NeedChoice {
			_, err = fmt.Scanln(&choice)
			if err != nil {
				os.Exit(0)
			}
			choice = strings.TrimSpace(choice)
		}

		nextRoute, err := route.Handler(m.cfg, choice)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		currentStep = nextRoute
	}
	os.Exit(0)

}
