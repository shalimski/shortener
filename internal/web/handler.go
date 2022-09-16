package web

import (
	"net/http"

	"github.com/shalimski/shortener/internal/ports"
	"github.com/shalimski/shortener/pkg/logger"
)

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
	h.log.Info(ctx, "create")
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info(ctx, "find")
}
