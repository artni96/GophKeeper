package menu

import (
	"context"
	"fmt"
	"strconv"

	cardspb "github.com/artni96/GophKeeper/api/proto/cards"
	"github.com/artni96/GophKeeper/internal/client/utils"
	"github.com/artni96/GophKeeper/internal/client/validators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (m *Menu) initCardMenu() {
	m.routes[cardList] = StepInfo{
		Logic: func(ctx context.Context) error {
			entityList := m.app.CardService.GetList()
			fmt.Println()
			fmt.Println("Number		Title		Description")
			for _, entity := range entityList {
				fmt.Printf("%d       %s       %s\n", entity.Number, entity.Title, entity.Description)
			}
			fmt.Println()
			fmt.Println("1. Get card by number")
			fmt.Println("2. Create new card")
			fmt.Println("3. Back")
			fmt.Println("0. Exit")
			fmt.Printf("Choose option: ")
			m.needInput = true

			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			switch choice {
			case "1":
				return getCardAskNumber, nil
			case "2":
				return createCard, nil
			case "3":
				return main, nil
			case "0":
				return exit, nil
			default:
				fmt.Println("Invalid choice")
				return cardList, nil
			}
		},
		NextSteps: []int{createCard, getCardAskNumber, main, exit},
	}

	m.routes[getCardAskNumber] = StepInfo{
		Logic: func(ctx context.Context) error {
			fmt.Printf("Enter card number: ")
			_, err := fmt.Scanln(&m.currentEntityNumber)
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}
			m.needInput = false
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return getCard, nil

		},
		NextSteps: []int{getCard},
	}

	m.routes[getCard] = StepInfo{
		Logic: func(ctx context.Context) error {
			entity, err := m.app.CardService.Get(m.currentEntityNumber)
			if err != nil {
				return err
			}
			fmt.Println()
			fmt.Printf("Number: %d\n", entity.Number)
			fmt.Printf("Title: %s\n", entity.Title)
			fmt.Printf("Description: %s\n", entity.Description)
			fmt.Printf("PAN: %d\n", entity.PAN)
			fmt.Printf("Holder: %s\n", entity.Holder)
			fmt.Printf("Expiry Date: %s\n", entity.ExpiryDate)
			fmt.Printf("CVV: %s\n", entity.CVV)
			fmt.Printf("PIN: %s\n", entity.PIN)
			fmt.Printf("Bank: %s\n", entity.Bank)
			fmt.Printf("Brand: %s\n", entity.Brand)
			fmt.Printf("Created at: %s\n", entity.CreatedAt)
			fmt.Printf("Updated at: %s\n", entity.UpdatedAt)
			fmt.Println()
			fmt.Println("1. Update")
			fmt.Println("2. Delete")
			fmt.Println("3. Back to menu")
			fmt.Printf("Choose option: ")

			m.needInput = true
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			switch choice {
			case "1":
				return updateCard, nil
			case "2":
				return deleteCard, nil
			case "3":
				return cardList, nil
			case "4":
				return main, nil
			default:
				fmt.Println("Invalid choice")
				return cardList, nil
			}

		},
		NextSteps: []int{updateCard, deleteCard, main},
	}

	m.routes[createCard] = StepInfo{
		Logic: func(ctx context.Context) error {
			m.needInput = false
			pbEntity := &cardspb.CardCreateRequest{}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter title: ")
				title, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid title")
					continue
				}
				if title == "" {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid title: field is required")
					continue
				}

				pbEntity.SetTitle(wrapperspb.String(title))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter description: ")
				description, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid description")
					continue
				}
				pbEntity.SetDescription(wrapperspb.String(description))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter PAN: ")
				strPAN, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PAN")
					continue
				}

				if strPAN == "" {
					break
				}

				uint64PAN, err := strconv.ParseUint(strPAN, 10, 64)
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PAN")
					continue
				}

				if err = validators.PANValidator(uint64PAN); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PAN")
					continue
				}
				pbEntity.SetPan(wrapperspb.UInt64(uint64PAN))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter card holder name: ")
				holder, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid card holder name")
					continue
				}
				pbEntity.SetHolder(wrapperspb.String(holder))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter Expiry Date: ")

				expiryDate, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid Expiry Date")
					continue
				}

				if expiryDate == "" {
					break
				}

				if err = validators.ExpiryDateValidator(expiryDate); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid Expiry Date")
					continue
				}
				pbEntity.SetExpiryDate(wrapperspb.String(expiryDate))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter CVV: ")
				cvv, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid CVV")
					continue
				}

				if cvv == "" {
					break
				}

				if err = validators.CVVValidator(cvv); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid CVV")
					continue
				}
				pbEntity.SetCvv(wrapperspb.String(cvv))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter PIN: ")
				pin, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PIN")
					continue
				}

				if pin == "" {
					break
				}

				if err = validators.PINValidator(pin); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PIN")
					continue

				}
				pbEntity.SetPin(wrapperspb.String(pin))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter bank name: ")
				bank, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid BANK")
				}
				pbEntity.SetBank(wrapperspb.String(bank))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter brand type (Visa, Mastercard, Mir...): ")
				brand, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid BANK")
					continue
				}
				pbEntity.SetBrand(wrapperspb.String(brand))
				break
			}

			fmt.Println()
			mdCtx := utils.PrepareMDContext(ctx, m.app.Cfg)

			_, err := m.app.CardService.Client.CreateCard(mdCtx, pbEntity)
			if err != nil {
				m.isFailed = true
				st, ok := status.FromError(err)
				if ok {
					if st.Code() == codes.AlreadyExists {
						fmt.Printf("Failed to create card record: %s\n", st.Message())
						return err
					}
				}
				fmt.Println("Failed to create card record")
				return err
			}
			fmt.Println("Card created successfully!")
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}

	m.routes[updateCard] = StepInfo{
		Logic: func(ctx context.Context) error {
			m.needInput = false
			pbEntity := &cardspb.CardUpdateRequest{}

			pbEntity.SetNumber(m.currentEntityNumber)

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter title: ")
				title, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid title")
					continue
				}

				pbEntity.SetTitle(wrapperspb.String(title))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter description: ")
				description, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid description")
					continue
				}
				pbEntity.SetDescription(wrapperspb.String(description))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter PAN: ")
				strPAN, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PAN")
					continue
				}

				if strPAN == "" {
					break
				}

				uint64PAN, err := strconv.ParseUint(strPAN, 10, 64)
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PAN")
					continue
				}

				if err = validators.PANValidator(uint64PAN); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PAN")
					continue
				}
				pbEntity.SetPan(wrapperspb.UInt64(uint64PAN))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter card holder name: ")
				holder, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid card holder name")
					continue
				}
				pbEntity.SetHolder(wrapperspb.String(holder))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter Expiry Date: ")

				expiryDate, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid Expiry Date")
					continue
				}

				if expiryDate == "" {
					break
				}

				if err = validators.ExpiryDateValidator(expiryDate); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid Expiry Date")
					continue
				}
				pbEntity.SetExpiryDate(wrapperspb.String(expiryDate))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter CVV: ")
				cvv, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid CVV")
					continue
				}

				if cvv == "" {
					break
				}

				if err = validators.CVVValidator(cvv); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid CVV")
					continue
				}
				pbEntity.SetCvv(wrapperspb.String(cvv))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter PIN: ")
				pin, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PIN")
					continue
				}

				if pin == "" {
					break
				}

				if err = validators.PINValidator(pin); err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid PIN")
					continue

				}
				pbEntity.SetPin(wrapperspb.String(pin))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter bank name: ")
				bank, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid BANK")
				}
				pbEntity.SetBank(wrapperspb.String(bank))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter brand type (Visa, Mastercard, Mir...): ")
				brand, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
					}
					fmt.Println("Invalid BANK")
					continue
				}
				pbEntity.SetBrand(wrapperspb.String(brand))
				break
			}

			err := m.confirmAction()
			if err != nil {
				return err
			}

			fmt.Println()
			mdCtx := utils.PrepareMDContext(ctx, m.app.Cfg)

			_, err = m.app.CardService.Client.UpdateCard(mdCtx, pbEntity)
			if err != nil {
				m.isFailed = true
				st, ok := status.FromError(err)
				if ok {
					if st.Code() == codes.AlreadyExists {
						fmt.Printf("Failed to create card: %s\n", st.Message())
						return err
					}
				}
				fmt.Println("Failed to update card")
				return err
			}
			fmt.Println("Card updated successfully!")
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}

	m.routes[deleteCard] = StepInfo{
		Logic: func(ctx context.Context) error {
			err := m.confirmAction()
			if err != nil {
				return err
			}

			m.needInput = false
			req := &cardspb.CardDeleteRequest{}
			req.SetNumber(m.currentEntityNumber)

			mdCtx := utils.PrepareMDContext(ctx, m.app.Cfg)
			_, err = m.app.CardService.Client.DeleteCard(mdCtx, req)
			if err != nil {
				m.isFailed = true
				fmt.Println("Failed to delete card record")
				return err
			}

			fmt.Println()
			fmt.Println("Card deleted successfully!")
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}
}
