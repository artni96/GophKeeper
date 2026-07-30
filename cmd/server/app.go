package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artni96/GophKeeper/internal/config"
	"github.com/artni96/GophKeeper/internal/constants"
	"github.com/artni96/GophKeeper/internal/facade"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// run initializes and starts the app.
func run(ctx context.Context, cfg *config.Config) error {
	eg := new(errgroup.Group)
	isClosedChan := make(chan struct{})

	sig, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	defer stop()

	app := facade.NewApp(eg, cfg)

	err := app.Launch(ctx)
	if err != nil {
		return err
	}

	<-sig.Done()
	app.Logger.Info("shutting down signal received")

	gfCtx, gfCancel := context.WithTimeout(ctx, cfg.GPeriod)
	defer gfCancel()

	app.Stop(gfCtx, gfCancel, isClosedChan)

	eg.Go(func() error {
		select {
		case <-isClosedChan:
			app.Logger.Info("the app stopped gracefully")

		case <-gfCtx.Done():
			app.Logger.Info("graceful period expired")

			app.Logger.Info(fmt.Sprintf("the app will be stopped forcefully in %v sec", app.Cfg.FPeriod))
			ffCtx, ffCancel := context.WithTimeout(ctx, app.Cfg.FPeriod)
			defer ffCancel()

			go ffCountdown(ffCtx, app.Logger, isClosedChan)

			select {
			case <-ffCtx.Done():
				app.Logger.Info("the app stopped forcefully")
				os.Exit(0)
			case _, ok := <-isClosedChan:
				if !ok {
					app.Logger.Info("the app stopped during the forceful period")
					ffCancel()
				}
			}

		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

// ffCountdown counts down left time of forceful shutdown.
func ffCountdown(ctx context.Context, logger *zap.Logger, gsChan <-chan struct{}) {
	deadline, _ := ctx.Deadline()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-gsChan:
			return
		case <-ticker.C:

			timeLeft := int(math.Ceil(time.Until(deadline).Seconds()))
			if timeLeft == 0 {
				return
			}
			if timeLeft == 15 || timeLeft == 10 || (timeLeft <= 5 && timeLeft > 0) {
				logger.Info(fmt.Sprintf("the app will be stopped forcefully in %d sec", timeLeft), zap.String("time", time.Now().Format(constants.TimeStamp)))
			}
		}
	}
}
