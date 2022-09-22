package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shalimski/shortener/config"
	"github.com/shalimski/shortener/internal/adapters/cache"

	"github.com/shalimski/shortener/internal/adapters/repository/urlrepo"
	"github.com/shalimski/shortener/internal/adapters/urlgenerator/generator"
	"github.com/shalimski/shortener/internal/services"
	"github.com/shalimski/shortener/internal/web"
	"github.com/shalimski/shortener/pkg/coordinator"
	"github.com/shalimski/shortener/pkg/httpserver"
	"github.com/shalimski/shortener/pkg/logger"
	"github.com/shalimski/shortener/pkg/mongodb"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	// Logger init

	ctx := context.Background()
	log := logger.NewLogger()
	log.With(zap.String("node", cfg.Node.Name))
	log.Info(ctx, "starting app...")

	// Database init
	mongoClient, err := mongodb.NewClient(cfg)
	if err != nil {
		log.Fatal("failed to connect MongoDB", zap.Error(err))
		return
	}
	db := urlrepo.NewURLRepo(mongoClient.Database(cfg.Mongo.Database))
	log.Info(ctx, "MongoDB initialized")

	// Coordinator for distributed counter
	counter, err := coordinator.NewCoordinator(cfg.App.EtcdEndpoints)
	defer counter.Shutdown()
	if err != nil {
		log.Fatal("failed to start distributed counter", zap.Error(err))
		return
	}
	log.Info(ctx, "distributed counter initialized")

	// Generator
	urlgen, err := generator.NewUrlGenerator(counter)
	if err != nil {
		log.Fatal("failed to init url generator", zap.Error(err))
		return
	}
	log.Info(ctx, "url generator initialized")

	redis := cache.NewCache(cfg)

	// Main service
	service := services.NewService(log, db, urlgen, redis)
	log.Info(ctx, "service initialized")

	h := web.NewHandler(service, log)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestID)
	r.Use(logger.Middleware(log))
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", h.Create)
		r.Get("/{shortURL}", h.Find)
		r.Delete("/{shortURL}", h.Delete)
	})

	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))
	log.Info(ctx, "http service started on port: "+cfg.HTTP.Port)

	log.Info(ctx, "-- Ready to accept connections --")

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
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(ctx, "failed to shutdown", zap.Error(err))
	}

	counter.Shutdown()
	redis.Shutdown(ctx)
}
