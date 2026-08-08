package menu

import "fmt"

func (m *Menu) initCardList() {
	m.routes[cardList] = StepInfo{
		Output: func(cfg *config) error {
			entityList := m.cfg.app.CardService.GetList()
			fmt.Println("Number	Title		Description")
			for _, entity := range entityList {
				fmt.Printf("%d       %s       %s\n", entity.Number, entity.Title, entity.Description)
			}
			fmt.Println()
			fmt.Println("1. Get card by number")
			fmt.Println("2. Create new card")
			fmt.Println("0. Exit")
			fmt.Printf("Choose option: ")
			m.cfg.NeedChoice = true

			return nil
		},
		Handler: func(cfg *config, choice string) (int, error) {
			switch choice {
			case "1":
				return getCard, nil
			case "2":
				return createCard, nil
			case "0":
				return exit, nil
			default:
				fmt.Println("Invalid choice")
				return cardList, nil
			}
		},
		NextSteps: []int{createCard, getCard, exit},
	}
}
