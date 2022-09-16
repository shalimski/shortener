package ports

import "github.com/shalimski/shortener/internal/domain"

type URLShortenerService interface {
	Create(longURL string) (shortURL string, err error)
	Find(shortURL string) (longURL string, err error)
}

type LinksRepository interface {
	Create(url domain.URL) error
	Find(shortURL string) (domain.URL, error)
}
