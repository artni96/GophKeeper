package menu

import (
	"context"
	"fmt"
)

func (m *Menu) initLoginList() {
	m.routes[listLogin] = StepInfo{
		Logic: func(ctx context.Context) error {
			entityList := m.app.LoginService.GetList()
			fmt.Println("Number	Title		Description")
			for _, entity := range entityList {
				fmt.Printf("%d       %s       %s\n", entity.Number, entity.Title, entity.Description)
			}
			fmt.Println()
			fmt.Println("1. Get card by number")
			fmt.Println("2. Create new card")
			fmt.Println("0. Exit")
			fmt.Printf("Choose option: ")
			m.NeedInput = true

			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			switch choice {
			case "1":
				return getLogin, nil
			case "2":
				return createLogin, nil
			case "0":
				return exit, nil
			default:
				fmt.Println("Invalid choice")
				return listLogin, nil
			}
		},
		NextSteps: []int{createCard, getCard, exit},
	}
}
