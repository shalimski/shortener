package memgen

import (
	"context"

	"github.com/shalimski/shortener/pkg/randomstring"
)

type urlGenerator struct {
	len int
}

func NewUrlGenerator(len int) *urlGenerator {
	return &urlGenerator{len: len}
}

func (u *urlGenerator) Next(ctx context.Context) string {
	return randomstring.New(u.len)
}