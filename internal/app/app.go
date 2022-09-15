package app

import (
	"context"

	"github.com/shalimski/shortener/config"
	"github.com/shalimski/shortener/pkg/logger"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	// logger init

	ctx := context.Background()
	log := logger.NewLogger()
	log.With(zap.String("node", cfg.Node.Name))
	log.Info(ctx, "starting app...")

	//
}
