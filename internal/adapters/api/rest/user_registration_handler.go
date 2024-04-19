package rest

import (
	"io"
	"net/http"

	"github.com/k0st1a/gophermart/internal/pkg/user"
	"github.com/rs/zerolog/log"
)

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

	rw.WriteHeader(http.StatusOK)
}
