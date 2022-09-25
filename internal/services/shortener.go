package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/shalimski/shortener/internal/domain"
	"github.com/shalimski/shortener/internal/ports"
	"github.com/shalimski/shortener/pkg/logger"
	"go.uber.org/zap"
)

// servise implements URLShortenerService
type service struct {
	log    *logger.Logger
	repo   ports.Repository
	urlgen ports.ShortURLGenerator
	cache  ports.Cacher
}

func NewService(log *logger.Logger, repo ports.Repository, urlgen ports.ShortURLGenerator, cache ports.Cacher) service {
	return service{
		log:    log,
		repo:   repo,
		urlgen: urlgen,
		cache:  cache,
	}
}

func (s service) Create(ctx context.Context, longURL string) (string, error) {
	s.log.Debug(ctx, "start Create method", zap.String("longURL", longURL))

	shortURL, err := s.urlgen.Next(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get next short url: %w", err)
	}

	s.log.Debug(ctx, "generated url", zap.String("shortURL", shortURL))

	url := domain.URL{
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	if err := s.repo.Create(ctx, url); err != nil {
		s.log.Error(ctx, "failed to create url", zap.Error(err))

		return "", domain.ErrFailedToCreate
	}

	if err := s.cache.Set(ctx, shortURL, longURL); err != nil {
		s.log.Error(ctx, "failed to set in cache", zap.Error(err))
	}

	return shortURL, nil
}

func (s service) Find(ctx context.Context, shortURL string) (longURL string, err error) {
	s.log.Debug(ctx, "start Find method", zap.String("shortURL", shortURL))

	if longURL, err = s.cache.Get(ctx, shortURL); err == nil {
		return longURL, nil
	}

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		s.log.Error(ctx, "failed to get in cache", zap.Error(err))
	}

	url, err := s.repo.Find(ctx, shortURL)
	if err != nil {
		return "", err
	}

	if err = s.cache.Set(ctx, shortURL, longURL); err != nil {
		s.log.Error(ctx, "failed to get in cache", zap.Error(err))
	}

	return url.LongURL, nil
}

func (s service) Delete(ctx context.Context, shortURL string) error {
	s.log.Debug(ctx, "start Delete method", zap.String("shortURL", shortURL))

	if err := s.cache.Del(ctx, shortURL); err != nil {
		s.log.Error(ctx, "failed to del in cache", zap.Error(err))
	}

	return s.repo.Delete(ctx, shortURL)
}
