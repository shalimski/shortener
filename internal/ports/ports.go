package ports

import (
	"context"

	"github.com/shalimski/shortener/internal/domain"
)

type ShortenerService interface {
	Create(ctx context.Context, longURL string) (shortURL string, err error)
	Find(ctx context.Context, shortURL string) (longURL string, err error)
	Delete(ctx context.Context, shortURL string) error
}

type Repository interface {
	Create(ctx context.Context, url domain.URL) error
	Find(ctx context.Context, shortURL string) (domain.URL, error)
	Delete(ctx context.Context, shortURL string) error
}

type ShortURLGenerator interface {
	Next(ctx context.Context) (string, error)
}

type Cacher interface {
	Set(ctx context.Context, shortURL string, longURL string) (err error)
	Get(ctx context.Context, shortURL string) (longURL string, err error)
	Del(ctx context.Context, shortURL string) (err error)
}
