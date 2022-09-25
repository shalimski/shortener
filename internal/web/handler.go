package web

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shalimski/shortener/internal/domain"
	"github.com/shalimski/shortener/internal/ports"
	"github.com/shalimski/shortener/pkg/logger"
	"github.com/shalimski/shortener/pkg/urlvalidator"
	"go.uber.org/zap"
)

const shortURLParam = "shortURL"

type Handler struct {
	log                 *logger.Logger
	urlShortenerService ports.ShortenerService
}

func NewHandler(service ports.ShortenerService, log *logger.Logger) *Handler {
	return &Handler{
		urlShortenerService: service,
		log:                 log,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info(ctx, "start create handler")

	// Reading body
	var data CreateURLDTO

	err := Decode(r, &data)
	defer r.Body.Close()

	if err != nil {
		h.log.Info(ctx, "failed to parse body")
		err = Respond(ctx, w, NewResponse(err.Error()), http.StatusBadRequest)

		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	// Validation
	if !urlvalidator.IsURL(data.LongURL) {
		h.log.Info(ctx, "invalid long url", zap.String("longURL", data.LongURL))
		err = Respond(ctx, w, NewResponse("invalid long url"), http.StatusBadRequest)

		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	// Create short link
	shortURL, err := h.urlShortenerService.Create(ctx, data.LongURL)
	if err != nil {
		h.log.Error(ctx, "failed to create url", zap.String("longURL", data.LongURL), zap.Error(err))
		err = Respond(ctx, w, NewResponse("failed to create url"), http.StatusInternalServerError)

		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}
	}

	err = Respond(ctx, w, ResponseCreateDTO{ShortURL: shortURL}, http.StatusOK)
	if err != nil {
		h.log.Error(ctx, "failed to respond", zap.Error(err))
	}
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info(ctx, "start find handler")

	shortURL := chi.URLParam(r, shortURLParam)

	// Validation
	if !urlvalidator.IsShortURLSuffix(shortURL) {
		h.log.Info(ctx, "invalid short url", zap.String("shortURL", shortURL))

		err := Respond(ctx, w, NewResponse("invalid short url"), http.StatusBadRequest)
		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	longURL, err := h.urlShortenerService.Find(ctx, shortURL)

	if err != nil && errors.Is(err, domain.ErrNotFound) {
		err = Respond(ctx, w, NewResponse("short url not found"), http.StatusNotFound)
		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	if err != nil {
		h.log.Info(ctx, "failed to find", zap.String("shortURL", shortURL), zap.String("error", err.Error()))
		err = Respond(ctx, w, NewResponse("failed to find"), http.StatusInternalServerError)

		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info(ctx, "start delete handler")

	shortURL := chi.URLParam(r, shortURLParam)

	// Validation
	if !urlvalidator.IsShortURLSuffix(shortURL) {
		h.log.Info(ctx, "invalid short url", zap.String("shortURL", shortURL))

		err := Respond(ctx, w, NewResponse("invalid short url"), http.StatusBadRequest)
		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	err := h.urlShortenerService.Delete(ctx, shortURL)

	if err != nil && errors.Is(err, domain.ErrNotFound) {
		err = Respond(ctx, w, NewResponse("short url not found"), http.StatusNotFound)

		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	if err != nil {
		h.log.Info(ctx, "failed to delete", zap.String("shortURL", shortURL), zap.String("error", err.Error()))

		err = Respond(ctx, w, NewResponse("failed to delete"), http.StatusInternalServerError)

		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	err = Respond(ctx, w, NewResponse("url deleted"), http.StatusOK)
	if err != nil {
		h.log.Error(ctx, "failed to respond", zap.Error(err))
	}
}
