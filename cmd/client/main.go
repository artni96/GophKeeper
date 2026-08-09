package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	eg.Go(func() error {
		clientMenu.Run(ctx)
		return nil
	})

	sig, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	defer stop()
	<-sig.Done()
	if err = eg.Wait(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	os.Exit(0)
}
