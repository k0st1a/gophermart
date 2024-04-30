package rest

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/k0st1a/gophermart/internal/adapters/api/rest/models"
	"github.com/k0st1a/gophermart/internal/pkg/auth"
	"github.com/k0st1a/gophermart/internal/pkg/order"
	"github.com/k0st1a/gophermart/internal/pkg/user"
	"github.com/rs/zerolog/log"
)

type handler struct {
	auth  auth.UserAuthentication
	user  user.Managment
	order order.OrderManagment
}

func NewHandler(auth auth.UserAuthentication, user user.Managment, order order.OrderManagment) *handler {
	return &handler{
		auth:  auth,
		user:  user,
		order: order,
	}
}

func (h *handler) register(rw http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	ur, err := models.DeserializeRegister(data)
	if err != nil {
		log.Error().Err(err).Msg("user registration deserialize error")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	passwordHash, err := h.auth.GeneratePasswordHash(ur.Password)
	if err != nil {
		log.Error().Err(err).Msg("error of generate password hash")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := h.user.Create(r.Context(), ur.Login, passwordHash)
	if err != nil {
		if errors.Is(err, user.ErrLoginAlreadyBusy) {
			rw.WriteHeader(http.StatusConflict)
			return
		}

		log.Error().Err(err).Msg("error of create user")
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

func (h *handler) login(rw http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	ul, err := models.DeserializeLogin(data)
	if err != nil {
		log.Error().Err(err).Msg("user login deserialize error")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	userID, password, err := h.user.GetIDAndPassword(r.Context(), ul.Login)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			rw.WriteHeader(http.StatusConflict)
			return
		}

		log.Error().Err(err).Msg("error of get user id and password")
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

func (h *handler) createOrder(rw http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("body read error")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumber := string(data)
	err = goluhn.Validate(orderNumber)
	if err != nil {
		http.Error(rw, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	orderID, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("order number parsing error")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.order.CreateOrder(r.Context(), userID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrOrderAlreadyUploadedByThisUser):
			rw.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, order.ErrOrderAlreadyUploadedByAnotherUser):
			rw.WriteHeader(http.StatusConflict)
			return
		default:
			log.Error().Err(err).Msg("create order error")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusAccepted)
}

func (h *handler) getOrders(rw http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.order.GetOrders(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("error of get orders")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("orders:%+v", orders)

	var modelOrders models.Orders
	for _, o := range orders {
		modelOrders = append(modelOrders, models.Order(o))
	}
	log.Printf("modelOrders:%+v", modelOrders)

	if len(orders) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := models.SerializeOrders(&modelOrders)
	if err != nil {
		log.Error().Err(err).Msg("error of serialize orders")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("error of write orders")
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *handler) getBalance(rw http.ResponseWriter, r *http.Request) {
}

func (h *handler) createWithdraw(rw http.ResponseWriter, r *http.Request) {
}

func (h *handler) getWithdrawals(rw http.ResponseWriter, r *http.Request) {
}
