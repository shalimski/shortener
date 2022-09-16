package services

import (
	"github.com/shalimski/shortener/internal/ports"
	"github.com/shalimski/shortener/pkg/logger"
)

// servise implements URLShortenerService
type service struct {
	log  *logger.Logger
	repo *ports.LinksRepository
}

func NewService(log *logger.Logger, repo *ports.LinksRepository) service {
	return service{log: log, repo: repo}
}

func (s service) Create(longURL string) (shortURL string, err error) {
	return "", nil
}

func (s service) Find(shortURL string) (longURL string, err error) {
	return "", nil
}
