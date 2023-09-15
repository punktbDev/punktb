package main

import (
	"gitlab.com/freelance/punkt-b/backend/config"
	"gitlab.com/freelance/punkt-b/backend/internal/server"
	"gitlab.com/freelance/punkt-b/backend/pkg/log"
	"go.uber.org/zap"
)

func main() {
	cfg := config.GetConfig()
	log.Verbosity(cfg.LogLevel)
	srv, err := server.NewServer(cfg)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	if err = srv.Run(); err != nil {
		zap.L().Fatal(err.Error())
	}
}
