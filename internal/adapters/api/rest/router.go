package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k0st1a/gophermart/internal/pkg/auth"
)

func BuildRoute(h *handler, a auth.UserAuthentication) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route(`/api/user`, func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post(`/register`, h.register)
			r.Post(`/login`, h.login)
		})
		r.Group(func(r chi.Router) {
			r.Use(authenticate(a))
			r.Post(`/orders`, h.createOrder)
			r.Get(`/orsers`, h.getOrders)
			r.Get(`/balance`, h.getBalance)
			r.Post(`/balance/withdraw`, h.createWithdraw)
			r.Get(`/withdrawals`, h.getWithdrawals)
		})
	})

	return r
}
