package menu

import (
	"context"
	"fmt"

	textspb "github.com/artni96/GophKeeper/api/proto/texts"
	"github.com/artni96/GophKeeper/internal/client/constants"
	"github.com/artni96/GophKeeper/internal/client/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (m *Menu) initTextMenu() {
	m.routes[textMenu] = StepInfo{
		Logic: func(ctx context.Context) error {
			entityList := m.app.TextService.GetList()
			fmt.Println()
			fmt.Println("Number		Title		Description")
			for _, entity := range entityList {
				fmt.Printf("%d       %s       %s\n", entity.Number, entity.Title, entity.Description)
			}
			fmt.Println()
			fmt.Println("1. Get text by number")
			fmt.Println("2. Create new text record")
			fmt.Println("3. Back")
			fmt.Println("0. Exit")
			fmt.Printf("Choose option: ")
			m.needInput = true

			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			switch choice {
			case "1":
				return getTextAskNumber, nil
			case "2":
				return createText, nil
			case "3":
				return main, nil
			case "0":
				return exit, nil
			default:
				fmt.Println("Invalid choice")
				return textMenu, nil
			}
		},
		NextSteps: []int{createText, getTextAskNumber, main, exit},
	}

	m.routes[getTextAskNumber] = StepInfo{
		Logic: func(ctx context.Context) error {
			fmt.Printf("Enter Text record number: ")
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
			return getText, nil

		},
		NextSteps: []int{getText},
	}

	m.routes[getText] = StepInfo{
		Logic: func(ctx context.Context) error {
			entity, err := m.app.TextService.Get(m.currentEntityNumber)
			if err != nil {
				m.needInput = false
				return err
			}
			fmt.Println()
			fmt.Printf("Number: %d\n", entity.Number)
			fmt.Printf("Title: %s\n", entity.Title)
			fmt.Printf("Description: %s\n", entity.Description)
			fmt.Printf("Text: %s\n", entity.Text)
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
				return updateText, nil
			case "2":
				return deleteText, nil
			case "3":
				return textMenu, nil
			case "4":
				return main, nil
			default:
				fmt.Println("Invalid choice")
				return textMenu, nil
			}
		},
		NextSteps: []int{updateText, deleteText, textMenu},
	}

	m.routes[createText] = StepInfo{
		Logic: func(ctx context.Context) error {
			m.needInput = false
			pbEntity := &textspb.TextCreateRequest{}

			activeKey := m.app.State.ActiveKeyID

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter title: ")
				title, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid title")
					continue
				}
				if title == "" || title == " " {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						fmt.Println("Title cannot be empty")
						return constants.ErrInvalidInput
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
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid description")
					continue
				}
				pbEntity.SetDescription(wrapperspb.String(description))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter text: ")
				text, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid text")
					continue
				}

				encryptedValue, nonce, err := m.app.EncryptField(text)

				if err != nil {
					continue
				}

				pbEntity.SetTextValue(wrapperspb.Bytes(encryptedValue))
				pbEntity.SetTextNonce(wrapperspb.Bytes(nonce))
				pbEntity.SetTextKeyId(wrapperspb.UInt64(activeKey))
				break
			}

			fmt.Println()

			mdCtx := utils.PrepareMDContext(ctx, m.app.State.Token)

			_, err := m.app.TextService.Client.CreateText(mdCtx, pbEntity)
			if err != nil {
				m.isFailed = true
				st, ok := status.FromError(err)
				if ok {
					if st.Code() == codes.AlreadyExists {
						fmt.Printf("Failed to create text record: %s\n", st.Message())
						return err
					}
				}
				fmt.Println("Failed to create text record")
				fmt.Println()
				return err
			}
			fmt.Println("Text record created successfully!")
			fmt.Println()
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}

	m.routes[updateText] = StepInfo{
		Logic: func(ctx context.Context) error {
			m.needInput = false
			pbEntity := &textspb.TextUpdateRequest{}

			activeKey := m.app.State.ActiveKeyID

			pbEntity.SetNumber(m.currentEntityNumber)

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter title: ")
				title, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid title")
					continue
				}

				if title == "" {
					break
				}

				if title == " " {
					fmt.Println("Title cannot be empty")
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
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid description")
					continue
				}
				pbEntity.SetDescription(wrapperspb.String(description))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter text: ")
				text, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid text")
					continue
				}

				encryptedValue, nonce, err := m.app.EncryptField(text)

				if err != nil {
					continue
				}

				pbEntity.SetTextValue(wrapperspb.Bytes(encryptedValue))
				pbEntity.SetTextNonce(wrapperspb.Bytes(nonce))
				pbEntity.SetTextKeyId(wrapperspb.UInt64(activeKey))
				break
			}
			err := m.confirmAction()
			if err != nil {
				return err
			}

			fmt.Println()

			mdCtx := utils.PrepareMDContext(ctx, m.app.State.Token)

			_, err = m.app.TextService.Client.UpdateText(mdCtx, pbEntity)
			if err != nil {
				m.isFailed = true
				st, ok := status.FromError(err)
				if ok {
					if st.Code() == codes.AlreadyExists {
						fmt.Printf("Failed to update text record: %s\n", st.Message())
						fmt.Println()
						return err
					}
				}
				fmt.Println("Failed to update text record")
				fmt.Println()
				return err
			}
			fmt.Println("Text record updated successfully!")
			fmt.Println()
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}

	m.routes[deleteText] = StepInfo{
		Logic: func(ctx context.Context) error {
			err := m.confirmAction()
			if err != nil {
				return err
			}

			m.needInput = false
			req := &textspb.TextDeleteRequest{}
			req.SetNumber(m.currentEntityNumber)

			mdCtx := utils.PrepareMDContext(ctx, m.app.State.Token)
			_, err = m.app.TextService.Client.DeleteText(mdCtx, req)
			if err != nil {
				m.isFailed = true
				fmt.Println("Failed to delete text record")
				fmt.Println()
				return err
			}

			fmt.Println("Text record deleted successfully!")
			fmt.Println()
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}
}
