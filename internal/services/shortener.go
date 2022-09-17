package services

import (
	"context"

	"github.com/shalimski/shortener/internal/domain"
	"github.com/shalimski/shortener/internal/ports"
	"github.com/shalimski/shortener/pkg/logger"
	"go.uber.org/zap"
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

func (s service) Create(ctx context.Context, longURL string) (string, error) {
	s.log.Debug(ctx, "start Create method", zap.String("longURL", longURL))

	shortURL := s.urlgen.Next(ctx)
	s.log.Debug(ctx, "generated url", zap.String("shortURL", shortURL))

	url := domain.URL{
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	if err := s.repo.Create(url); err != nil {
		s.log.Error(ctx, "failed to create url", zap.Error(err))
		return "", ErrFailedToCreate
	}

	// TODO set in cache

	return shortURL, nil
}

func (s service) Find(ctx context.Context, shortURL string) (longURL string, err error) {
	s.log.Debug(ctx, "start Find method", zap.String("shortURL", shortURL))

	// TODO try in cache

	url, err := s.repo.Find(shortURL)
	if err != nil {
		return "", err
	}

	// TODO set in cache

	return url.LongURL, nil
}

func (s service) Delete(ctx context.Context, shortURL string) error {
	s.log.Debug(ctx, "start Delete method", zap.String("shortURL", shortURL))

	// TODO try in cache

	return s.repo.Delete(shortURL)
}
