package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/bootstrap"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	app := bootstrap.New()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		app.Shutdown(ctx)
		zap.L().Info("Received shutdown signal", zap.String("signal", sig.String()))
		os.Exit(0)
	}()

	app.Run()
}
