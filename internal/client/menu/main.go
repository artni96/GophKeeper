package menu

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"strings"

	"github.com/artni96/GophKeeper/internal/client/app"
)

const (
	initial = iota
	register
	login
	main
	exit

	cardMenu
	createCard
	getCardAskNumber
	getCard
	updateCard
	deleteCard

	loginMenu
	createLogin
	getLoginAskNumber
	getLogin
	updateLogin
	deleteLogin

	textMenu
	createText
	getTextAskNumber
	getText
	updateText
	deleteText
)

type StepInfo struct {
	Logic          func(ctx context.Context) error
	HandleNextStep func(choice string) (int, error)
	NextSteps      []int
}

type Menu struct {
	reader              *bufio.Reader
	routes              map[int]StepInfo
	app                 *app.App
	needInput           bool
	currentEntityNumber uint64
	isFailed            bool
}

func InitMenu(app *app.App) *Menu {
	return &Menu{
		reader: bufio.NewReader(os.Stdin),
		routes: make(map[int]StepInfo),
		app:    app,
	}
}

func (m *Menu) initSteps() {

	m.routes = make(map[int]StepInfo)

	m.initUserMenu()
	m.initCardMenu()
	m.initLoginMenu()
	m.initTextMenu()

	m.routes[initial] = StepInfo{
		Logic: func(ctx context.Context) error {

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
		NextSteps: []int{register, login, exit},
	}

	m.routes[main] = StepInfo{
		Logic: func(ctx context.Context) error {
			fmt.Println("1. Cards")
			fmt.Println("2. Logins")
			fmt.Println("3. Texts")
			fmt.Println("0. Logout & Exit")
			fmt.Print("Choose option: ")
			m.needInput = true
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			switch choice {
			case "1":
				return cardMenu, nil
			case "2":
				return loginMenu, nil
			case "3":
				return textMenu, nil
			case "0":
				m.app.State.IsOnline = false
				m.app.State.Token = ""
				fmt.Println("Logged out. Goodbye!")
				return exit, nil
			default:
				fmt.Println("Invalid choice")
				return main, nil
			}
		},
		NextSteps: []int{main, exit},
	}
}

func (m *Menu) confirmAction() error {
	fmt.Printf("Are you sure? (Y/n): ")
	choice, err := m.app.ReadLine()
	if err != nil {
		return err
	}

	switch strings.ToLower(choice) {
	case "y":
		return nil
	case "n":
		m.isFailed = true
		return fmt.Errorf("cancelled")
	default:
		m.isFailed = true
		return fmt.Errorf("cancelled")
	}
}

func (m *Menu) Run(ctx context.Context) {
	m.initSteps()
	currentStep := initial
	failsCounter := 0

	for {
		if m.app.IsBeingUpdated {
			continue
		}
		for currentStep != exit {
			route := m.routes[currentStep]

			err := route.Logic(ctx)
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
				} else if m.isFailed {
					m.isFailed = false
				} else {
					continue
				}

			}

			var choice string
			if m.needInput {

				choice, err = m.app.ReadLine()
				if err != nil {
					os.Exit(0)
				}
				choice = strings.TrimSpace(choice)

			}

			nextRoute, err := route.HandleNextStep(choice)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			currentStep = nextRoute
		}
		os.Exit(0)
	}

}
