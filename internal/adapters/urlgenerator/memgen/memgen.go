// Simplest implementation of a new random shortURL
package memgen

import (
	"context"

	"github.com/shalimski/shortener/pkg/randomstring"
)

type urlGenerator struct {
	length int
}

func NewURLGenerator(length int) *urlGenerator {
	return &urlGenerator{length: length}
}

func (u *urlGenerator) Next(ctx context.Context) (string, error) {
	return randomstring.New(u.length), nil
}
