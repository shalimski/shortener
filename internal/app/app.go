package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shalimski/shortener/config"
	memdb "github.com/shalimski/shortener/internal/adapters/repository/mem"
	"github.com/shalimski/shortener/internal/adapters/urlgenerator/memgen"
	"github.com/shalimski/shortener/internal/services"
	"github.com/shalimski/shortener/internal/web"
	"github.com/shalimski/shortener/pkg/httpserver"
	"github.com/shalimski/shortener/pkg/logger"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	// Logger init

	ctx := context.Background()
	log := logger.NewLogger()
	log.With(zap.String("node", cfg.Node.Name))
	log.Info(ctx, "starting app...")

	// Database init
	db := memdb.New()
	log.Info(ctx, "db init")
	//
	urlgen := memgen.NewUrlGenerator(cfg.App.ShortURLLength)
	log.Info(ctx, "url generator init")

	service := services.NewService(log, db, urlgen)
	log.Info(ctx, "service init")

	h := web.NewHandler(service, log)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestID)
	r.Use(logger.Middleware(log))
	r.Post("/new", h.Create)
	r.Get("/{shortURL}", h.Find)

	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))
	log.Info(ctx, "http service started on port: "+cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info(ctx, "signal: "+s.String())
	case err := <-httpServer.Notify():
		log.Error(ctx, "httpServer was stopped", zap.Error(err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		log.Error(ctx, "failed to shutdown", zap.Error(err))
	}
}
