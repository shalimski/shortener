package memdb

import (
	"context"
	"sync"

	"github.com/shalimski/shortener/internal/domain"
	"github.com/shalimski/shortener/internal/ports"
)

// basic realization for storage
type memdb struct {
	mu sync.RWMutex
	db map[string]string
}

func New() ports.Repository {
	return &memdb{db: make(map[string]string)}
}

func (m *memdb) Create(ctx context.Context, url domain.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.db[url.ShortURL] = url.LongURL

	return nil
}

func (m *memdb) Find(ctx context.Context, shortURL string) (domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	longURL, ok := m.db[shortURL]

	if !ok {
		return domain.URL{}, domain.ErrNotFound
	}

	u := domain.URL{
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	return u, nil
}

func (m *memdb) Delete(ctx context.Context, shortURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.db, shortURL)

	return nil
}
