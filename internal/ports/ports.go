package ports

import (
	"context"

	"github.com/shalimski/shortener/internal/domain"
)

type URLShortenerService interface {
	Create(ctx context.Context, longURL string) (shortURL string, err error)
	Find(ctx context.Context, shortURL string) (longURL string, err error)
	Delete(ctx context.Context, shortURL string) error
}

type LinksRepository interface {
	Create(ctx context.Context, url domain.URL) error
	Find(ctx context.Context, shortURL string) (domain.URL, error)
	Delete(ctx context.Context, shortURL string) error
}

type ShortURLGenerator interface {
	Next(ctx context.Context) (string, error)
}
