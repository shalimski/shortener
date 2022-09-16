package memdb

import (
	"fmt"
	"sync"

	"github.com/shalimski/shortener/internal/domain"
)

type memdb struct {
	mu sync.RWMutex
	db map[string]string
}

func New() *memdb {
	return &memdb{db: make(map[string]string)}
}

func (m *memdb) Create(url domain.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.db[url.ShortURL] = url.LongURL
	return nil
}

func (m *memdb) Find(shortURL string) (domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	longURL, ok := m.db[shortURL]
	if !ok {
		return domain.URL{}, fmt.Errorf("url not found")
	}

	u := domain.URL{
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	return u, nil
}

func (m *memdb) Delete(shortURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.db, shortURL)
	return nil
}