package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/utils/logger"

	"github.com/Format-C-eft/universal-proxy/internal/bootstrap"
)

func main() {
	ctx := context.Background()

	if config.VersionFlag {
		fmt.Printf("Name - '%s'\n", config.GetVersion().Name)
		fmt.Printf("Branch - '%s'\n", config.GetVersion().Branch)
		fmt.Printf("Commit hash - '%s'\n", config.GetVersion().CommitHash)
		fmt.Printf("Time build - '%s'\n", config.GetVersion().TimeBuild)
		return
	}

	servers, err := bootstrap.InitializeServers(ctx)
	if err != nil {
		logger.ErrorKV(ctx, "can`t init servers", "err", err)
		return
	}

	workers, err := bootstrap.InitializeWorker(ctx)
	if err != nil {
		logger.ErrorKV(ctx, "can`t init workers", "err", err)
		return
	}

	servers.Run(ctx)
	workers.Run(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.InfoF(ctx, "signal.Notify: %v", v)
	case done := <-ctx.Done():
		logger.InfoF(ctx, "ctx.Done: %v", done)
	}

	logger.Info(ctx, "Shutting down server...")

	servers.Stop(ctx)
	workers.Stop(ctx)
}
