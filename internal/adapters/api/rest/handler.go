package rest

import (
	"io"
	"net/http"
	"strconv"

	"github.com/k0st1a/gophermart/internal/pkg/user"
	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type handler struct {
	storage ports.UserStorage
}

func newHandler(s ports.UserStorage) *handler {
	return &handler{
		storage: s,
	}
}

func (h *handler) userRegistrationHandler(rw http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("uri", r.RequestURI).
		Str("method", r.Method).
		Msg("")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	ur, err := user.DeserializeRegister(b)
	if err != nil {
		log.Error().Err(err).Msg("user registration deserialize error")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	err = ur.Validate()
	if err != nil {
		log.Error().Err(err).Msg("user registration validation error")
		http.Error(rw, "validation error", http.StatusBadRequest)
		return
	}

	id, err := h.storage.CreateUser(r.Context(), ur.Login, ur.Password)
	if err != nil {
		log.Error().Err(err).Msg("user create error")
		http.Error(rw, "create user error", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Authorization", strconv.FormatInt(id, 10))
	rw.WriteHeader(http.StatusOK)
}
