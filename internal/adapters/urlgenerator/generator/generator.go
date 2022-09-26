// Distributed short URL generator,
// each node has an ints interval for generating links in order.
// Intervals of different nodes do not overlap, a coordinator is responsible for this
package generator

import (
	"context"
	"sync"

	"github.com/shalimski/shortener/internal/ports"
)

const (
	interval = 100000
	base     = 62
	chars    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type urlGenerator struct {
	mu          sync.Mutex
	currCounter int // current value of interval, Encode(current) will be next shortURL
	maxCounter  int // max value of interval, when will be reached, it need to request new interval

	counter Counter // distributed counter
}

type Counter interface {
	NextCounter(context.Context) (int, error)
}

func NewURLGenerator(counter Counter) (ports.ShortURLGenerator, error) {
	u := &urlGenerator{
		counter: counter,
	}

	ctx := context.Background()
	if err := u.setNextInterval(ctx); err != nil {
		return nil, err
	}

	return u, nil
}

// Next value of short URL
func (u *urlGenerator) Next(ctx context.Context) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.currCounter == u.maxCounter {
		if err := u.setNextInterval(ctx); err != nil {
			return "", err
		}
	}

	shortURL := Encode(u.currCounter)
	u.currCounter++

	return shortURL, nil
}

// setNextInterval get next value of distributed counter and set next interval based on it
func (u *urlGenerator) setNextInterval(ctx context.Context) error {
	next, err := u.counter.NextCounter(ctx)
	if err != nil {
		return err
	}

	u.currCounter = 1 + interval*(next-1)
	u.maxCounter = interval * next

	return nil
}

// Encode returns base62 representation of int
func Encode(num int) string {
	if num < 0 {
		return ""
	}

	result := make([]byte, 0)

	for num > 0 {
		curr := num % base
		num /= base

		result = append([]byte{chars[curr]}, result...)
	}

	return string(result)
}
