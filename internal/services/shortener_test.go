package services_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shalimski/shortener/internal/domain"
	mock "github.com/shalimski/shortener/internal/ports/mock"
	"github.com/shalimski/shortener/internal/services"
	"github.com/shalimski/shortener/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	log := logger.NewDebugLogger()

	url := domain.URL{
		ShortURL: "abcd",
		LongURL:  "http://github.com",
	}
	repo := mock.NewMockRepository(ctl)
	repo.EXPECT().Create(ctx, url).Return(nil)

	urlgen := mock.NewMockShortURLGenerator(ctl)
	urlgen.EXPECT().Next(ctx).Return(url.ShortURL, nil)

	cache := mock.NewMockCacher(ctl)
	cache.EXPECT().Set(ctx, url.ShortURL, url.LongURL).Return(nil)

	service := services.NewService(log, repo, urlgen, cache)
	shortURL, err := service.Create(ctx, url.LongURL)

	assert.NoError(t, err)
	assert.Equal(t, url.ShortURL, shortURL)
}

func TestFind(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	log := logger.NewDebugLogger()

	url := domain.URL{
		ShortURL: "abcd",
		LongURL:  "http://github.com",
	}
	urlgen := mock.NewMockShortURLGenerator(ctl)

	repo := mock.NewMockRepository(ctl)
	repo.EXPECT().Find(ctx, url.ShortURL).Return(url, nil)

	cache := mock.NewMockCacher(ctl)
	cache.EXPECT().Get(ctx, url.ShortURL).Return("", domain.ErrNotFound)
	cache.EXPECT().Set(ctx, url.ShortURL, url.LongURL).Return(nil)

	service := services.NewService(log, repo, urlgen, cache)
	longURL, err := service.Find(ctx, url.ShortURL)

	assert.NoError(t, err)
	assert.Equal(t, url.LongURL, longURL)
}

func TestDelete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	log := logger.NewDebugLogger()

	url := domain.URL{
		ShortURL: "abcd",
		LongURL:  "http://github.com",
	}
	urlgen := mock.NewMockShortURLGenerator(ctl)

	repo := mock.NewMockRepository(ctl)
	repo.EXPECT().Delete(ctx, url.ShortURL).Return(nil)

	cache := mock.NewMockCacher(ctl)
	cache.EXPECT().Del(ctx, url.ShortURL).Return(nil)

	service := services.NewService(log, repo, urlgen, cache)
	err := service.Delete(ctx, url.ShortURL)

	assert.NoError(t, err)
}
