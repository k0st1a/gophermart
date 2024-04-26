package rest

import (
	"errors"
	"io"
	"net/http"

	"github.com/k0st1a/gophermart/internal/pkg/auth"
	"github.com/k0st1a/gophermart/internal/pkg/user"
	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type handler struct {
	storage ports.UserStorage
	auth    auth.UserAuthenticator
}

func newHandler(s ports.UserStorage) *handler {
	return &handler{
		storage: s,
	}
}

func (h *handler) userRegistrationHandler(rw http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	ur, err := user.DeserializeRegister(b)
	if err != nil {
		log.Error().Err(err).Msg("user registration deserialize error")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ur.Validate()
	if err != nil {
		log.Error().Err(err).Msg("user registration validation error")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	ph, err := h.auth.GeneratePasswordHash(ur.Password)
	if err != nil {
		log.Error().Err(err).Msg("error of generate password hash")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := h.storage.CreateUser(r.Context(), ur.Login, ph)
	if err != nil {
		if errors.Is(err, ports.ErrLoginAlreadyBusy) {
			rw.WriteHeader(http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("user create error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	t, err := h.auth.GenerateToken(id)
	if err != nil {
		log.Error().Err(err).Msg("error of generate token")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Authorization", t)
	rw.WriteHeader(http.StatusOK)
}

func (h *handler) userAuthorizationHandler(rw http.ResponseWriter, r *http.Request) {
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

	ul, err := user.DeserializeLogin(b)
	if err != nil {
		log.Error().Err(err).Msg("user login deserialize error")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	err = ul.Validate()
	if err != nil {
		log.Error().Err(err).Msg("user login validation error")
		http.Error(rw, "validation error", http.StatusBadRequest)
		return
	}

	userID, password, err := h.storage.GetUserIDAndPassword(r.Context(), ul.Login)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			rw.WriteHeader(http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("error of get user password")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.auth.CheckPasswordHash(ul.Password, password)
	if err != nil {
		rw.WriteHeader(http.StatusConflict)
		return
	}

	t, err := h.auth.GenerateToken(userID)
	if err != nil {
		log.Error().Err(err).Msg("error of generate token")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Authorization", t)
	rw.WriteHeader(http.StatusOK)
}
