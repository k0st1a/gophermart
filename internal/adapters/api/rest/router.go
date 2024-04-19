package rest

import (
	"github.com/go-chi/chi/v5"
)

func newRouter() *chi.Mux {
	return chi.NewRouter()
}

func buildRoute(r *chi.Mux, h *handler) {
	r.Post("/api/user/register", h.userRegistrationHandler)
}
