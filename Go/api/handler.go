package api

import (
	"oauth/Go/common/config"
	"oauth/Go/common/routes"
	db "oauth/migrations/sqlc"

	"github.com/go-chi/chi/v5"
)

// Handler combines all route handlers
type Handler struct {
	config *config.Config
	db     db.Querier
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, querier db.Querier) *Handler {
	return &Handler{
		config: cfg,
		db:     querier,
	}
}

// Routes sets up all HTTP routes using chi router
func (h *Handler) Routes() *chi.Mux {
	router := routes.DefaultRouter()

	router.Post("/sso/request-otp", h.RequestOTP)
    router.Post("/sso/login", h.Login)

	return router
}