package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"stocks/internal/bootstrap"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.NewApp(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to initialize application \n")
		os.Exit(1)
	}

	logger := app.Logger()

	if err := app.Run(); err != nil {
		logger.Fatalf("❌ Application stopped with error: %v", err)
	}

	logger.Info("✅ Application exited successfully.")
}
