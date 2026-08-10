package menu

import (
	"context"
	"fmt"

	loginspb "github.com/artni96/GophKeeper/api/proto/logins"
	"github.com/artni96/GophKeeper/internal/client/utils"
	"github.com/artni96/GophKeeper/internal/server/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (m *Menu) initLoginMenu() {
	m.routes[loginMenu] = StepInfo{
		Logic: func(ctx context.Context) error {
			entityList := m.app.LoginService.GetList()
			fmt.Println()
			fmt.Println("Number		Title		Description")
			for _, entity := range entityList {
				fmt.Printf("%d       %s       %s\n", entity.Number, entity.Title, entity.Description)
			}
			fmt.Println()
			fmt.Println("1. Get login by number")
			fmt.Println("2. Create new login record")
			fmt.Println("3. Back")
			fmt.Println("0. Exit")
			fmt.Printf("Choose option: ")
			m.needInput = true

			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			switch choice {
			case "1":
				return getLoginAskNumber, nil
			case "2":
				return createLogin, nil
			case "3":
				return main, nil
			case "0":
				return exit, nil
			default:
				fmt.Println("Invalid choice")
				return loginMenu, nil
			}
		},
		NextSteps: []int{createLogin, getLoginAskNumber, main, exit},
	}

	m.routes[getLoginAskNumber] = StepInfo{
		Logic: func(ctx context.Context) error {
			fmt.Printf("Enter Login record number: ")
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
			return getLogin, nil

		},
		NextSteps: []int{getLogin},
	}

	m.routes[getLogin] = StepInfo{
		Logic: func(ctx context.Context) error {
			entity, err := m.app.LoginService.Get(m.currentEntityNumber)
			if err != nil {
				m.needInput = false
				return err
			}
			fmt.Println()
			fmt.Printf("Number: %d\n", entity.Number)
			fmt.Printf("Title: %s\n", entity.Title)
			fmt.Printf("Description: %s\n", entity.Description)
			fmt.Printf("Login: %s\n", entity.Login)
			fmt.Printf("Password: %s\n", entity.Password)
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
				return updateLogin, nil
			case "2":
				return deleteLogin, nil
			case "3":
				return loginMenu, nil
			case "4":
				return main, nil
			default:
				fmt.Println("Invalid choice")
				return loginMenu, nil
			}

		},
		NextSteps: []int{updateLogin, deleteLogin, main},
	}

	m.routes[createLogin] = StepInfo{
		Logic: func(ctx context.Context) error {
			m.needInput = false
			pbEntity := &loginspb.LoginCreateRequest{}

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
				fmt.Printf("Enter login: ")
				loginVal, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid login")
					continue
				}

				encryptedValue, nonce, err := m.app.EncryptField(loginVal)

				if err != nil {
					continue
				}

				pbEntity.SetLogin(wrapperspb.Bytes(encryptedValue))
				pbEntity.SetLoginNonce(wrapperspb.Bytes(nonce))
				pbEntity.SetLoginKeyId(wrapperspb.UInt64(activeKey))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter password: ")
				password, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid password")
					continue
				}

				encryptedValue, nonce, err := m.app.EncryptField(password)

				if err != nil {
					continue
				}

				pbEntity.SetPassword(wrapperspb.Bytes(encryptedValue))
				pbEntity.SetPasswordNonce(wrapperspb.Bytes(nonce))
				pbEntity.SetPasswordKeyId(wrapperspb.UInt64(activeKey))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter URL: ")
				url, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid URL")
					continue
				}
				pbEntity.SetUrl(wrapperspb.String(url))
				break
			}

			fmt.Println()

			mdCtx := utils.PrepareMDContext(ctx, m.app.State.Token)

			_, err := m.app.LoginService.Client.CreateLogin(mdCtx, pbEntity)
			if err != nil {
				m.isFailed = true
				st, ok := status.FromError(err)
				if ok {
					if st.Code() == codes.AlreadyExists {
						fmt.Printf("Failed to create login record: %s\n", st.Message())
						return err
					}
				}
				fmt.Println("Failed to create login record")
				fmt.Println()
				return err
			}
			fmt.Println("Login record created successfully!")
			fmt.Println()
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}

	m.routes[updateLogin] = StepInfo{
		Logic: func(ctx context.Context) error {
			m.needInput = false
			pbEntity := &loginspb.LoginUpdateRequest{}

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
				fmt.Printf("Enter login: ")
				loginVal, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid login")
					continue
				}

				encryptedValue, nonce, err := m.app.EncryptField(loginVal)

				if err != nil {
					continue
				}

				pbEntity.SetLogin(wrapperspb.Bytes(encryptedValue))
				pbEntity.SetLoginNonce(wrapperspb.Bytes(nonce))
				pbEntity.SetLoginKeyId(wrapperspb.UInt64(activeKey))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter password: ")
				password, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid password")
					continue
				}

				encryptedValue, nonce, err := m.app.EncryptField(password)

				if err != nil {
					continue
				}

				pbEntity.SetPassword(wrapperspb.Bytes(encryptedValue))
				pbEntity.SetPasswordNonce(wrapperspb.Bytes(nonce))
				pbEntity.SetPasswordKeyId(wrapperspb.UInt64(activeKey))
				break
			}

			for i := range m.app.Cfg.MaxAttempts {
				fmt.Printf("Enter URL: ")
				url, err := m.app.ReadLine()
				if err != nil {
					if i == m.app.Cfg.MaxAttempts-1 {
						m.isFailed = true
						return constants.ErrInvalidInput
					}
					fmt.Println("Invalid URL")
					continue
				}
				pbEntity.SetUrl(wrapperspb.String(url))
				break
			}

			pbEntity.SetNumber(m.currentEntityNumber)

			fmt.Println()

			mdCtx := utils.PrepareMDContext(ctx, m.app.State.Token)

			_, err := m.app.LoginService.Client.UpdateLogin(mdCtx, pbEntity)
			if err != nil {
				m.isFailed = true
				st, ok := status.FromError(err)
				if ok {
					if st.Code() == codes.AlreadyExists {
						fmt.Printf("Failed to update login record: %s\n", st.Message())
						fmt.Println()
						return err
					}
				}
				fmt.Println("Failed to update login record")
				fmt.Println()
				return err
			}
			fmt.Println("Login record updated successfully!")
			fmt.Println()
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}

	m.routes[deleteLogin] = StepInfo{
		Logic: func(ctx context.Context) error {
			err := m.confirmAction()
			if err != nil {
				return err
			}

			m.needInput = false
			req := &loginspb.LoginDeleteRequest{}
			req.SetNumber(m.currentEntityNumber)

			mdCtx := utils.PrepareMDContext(ctx, m.app.State.Token)
			_, err = m.app.LoginService.Client.DeleteLogin(mdCtx, req)
			if err != nil {
				m.isFailed = true
				fmt.Println("Failed to delete login record")
				fmt.Println()
				return err
			}

			fmt.Println("Login record deleted successfully!")
			fmt.Println()
			return nil
		},
		HandleNextStep: func(choice string) (int, error) {
			return main, nil
		},
		NextSteps: []int{main},
	}
}
