package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shalimski/shortener/internal/ports"
	"github.com/shalimski/shortener/pkg/logger"
	"github.com/shalimski/shortener/pkg/urlvalidator"
	"go.uber.org/zap"
)

const shortURLParam = "shorturl"

type Handler struct {
	log                 *logger.Logger
	urlShortenerService ports.URLShortenerService
}

func NewHandler(service ports.URLShortenerService, log *logger.Logger) *Handler {
	return &Handler{
		urlShortenerService: service,
		log:                 log,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Debug(ctx, "start create handler")

	// Reading body
	var data CreateUrlDTO

	err := Decode(r, &data)
	defer r.Body.Close()
	if err != nil {
		h.log.Info(ctx, "failed to parse body")
		err = Respond(ctx, w, NewResponseError(err.Error()), http.StatusBadRequest)
		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	// Validation
	if !urlvalidator.IsURL(data.LongURL) {
		h.log.Info(ctx, "invalid long url", zap.String("longURL", data.LongURL))
		err = Respond(ctx, w, NewResponseError("invalid long url"), http.StatusBadRequest)
		if err != nil {
			h.log.Error(ctx, "failed to respond", zap.Error(err))
		}

		return
	}

	// Create short link
	shortURL, err := h.urlShortenerService.Create(ctx, data.LongURL)
	if err != nil {
		h.log.Error(ctx, "failed to create url", zap.String("longURL", data.LongURL), zap.Error(err))
		err = Respond(ctx, w, NewResponseError("failed to create url"), http.StatusInternalServerError)
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
	h.log.Debug(ctx, "start find handler")

	shortURL := chi.URLParam(r, shortURLParam)

	_, err := h.urlShortenerService.Find(ctx, shortURL)
	if err != nil {
		// 404
	}

	// 301

	// Respond(ctx, w, NewResponseError(err.Error()), http.StatusOK)
	// longURL
}
