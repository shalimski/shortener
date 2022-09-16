package services

import (
	"context"

	"github.com/shalimski/shortener/internal/ports"
	"github.com/shalimski/shortener/pkg/logger"
)

// servise implements URLShortenerService
type service struct {
	log    *logger.Logger
	repo   ports.LinksRepository
	urlgen ports.ShortURLGenerator
}

func NewService(log *logger.Logger, repo ports.LinksRepository, urlgen ports.ShortURLGenerator) service {
	return service{
		log:    log,
		repo:   repo,
		urlgen: urlgen,
	}
}

func (s service) Create(ctx context.Context, longURL string) (shortURL string, err error) {
	s.log.Info(ctx, "create")
	shorturl := s.urlgen.Next(ctx)

	s.log.Info(ctx, "new short url: "+shorturl)
	return shorturl, nil
}

func (s service) Find(ctx context.Context, shortURL string) (longURL string, err error) {
	s.log.Info(ctx, "find")
	return "", nil
}
