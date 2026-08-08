package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	userspb "github.com/artni96/GophKeeper/api/proto/users"
	clientapp "github.com/artni96/GophKeeper/internal/client/app"
	"github.com/artni96/GophKeeper/internal/client/menu"
	"golang.org/x/sync/errgroup"
)

func main() {
	eg := errgroup.Group{}
	ctx := context.Background()
	app := clientapp.NewApp(&eg)

	err := app.InitServerConn(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer app.CloseServerConn()

	app.InitDependencies()

	clientMenu := menu.InitMenu(app)
	clientMenu.Run(ctx)

	//err = app.UploadData(ctx)
	//if err != nil {
	//	fmt.Println("failed to upload data")
	//	os.Exit(0)
	//}

	sig, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	defer stop()
	<-sig.Done()
	if err = eg.Wait(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	os.Exit(0)
}

func runGRPCClient(ctx context.Context, eg *errgroup.Group) {
	for {
		select {
		case <-ctx.Done():
			if err := eg.Wait(); err != nil {
				fmt.Println(err)
			}
			time.Sleep(1 * time.Second)
			fmt.Println("Client stopping...")
			return
		default:
			fmt.Println("Client running...")
			time.Sleep(1 * time.Second)
		}
	}
}

func auth(ctx context.Context, c userspb.UserServiceClient) error {
	var login string
	var password string
	fmt.Printf("enter login: ")
	fmt.Scan(&login)
	fmt.Printf("enter password: ")
	fmt.Scan(&password)

	req := &userspb.LoginRequest{}
	req.SetUsername(login)
	req.SetPassword(password)
	_, err := c.Login(ctx, req)
	if err != nil {
		return err
	}
	return nil

}
