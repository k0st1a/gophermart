package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/k0st1a/gophermart/internal/pkg/auth"
	"github.com/k0st1a/gophermart/internal/pkg/order"
	"github.com/k0st1a/gophermart/internal/pkg/user"
	"github.com/k0st1a/gophermart/internal/pkg/withdraw"
	"github.com/rs/zerolog/log"
)

type handler struct {
	auth     auth.UserAuthentication
	user     user.Managment
	order    order.Managment
	withdraw withdraw.Managment
}

func NewHandler(a auth.UserAuthentication, u user.Managment, o order.Managment, w withdraw.Managment) *handler {
	return &handler{
		auth:     a,
		user:     u,
		order:    o,
		withdraw: w,
	}
}

func (h *handler) register(rw http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ur Register
	err = json.Unmarshal(data, &ur)
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

	var ul Login
	err = json.Unmarshal(data, &ul)
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

	err = h.order.Create(r.Context(), userID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrAlreadyUploadedByThisUser):
			rw.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, order.ErrAlreadyUploadedByAnotherUser):
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

//nolint:dupl //similar to getWithdrawals
func (h *handler) getOrders(rw http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.order.List(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("error of get orders")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("orders:%+v", orders)

	if len(orders) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	modelOrders := make([]Order, len(orders))
	for i := 0; i < len(orders); i++ {
		modelOrders[i] = Order(orders[i])
	}
	log.Printf("modelOrders:%+v", modelOrders)

	data, err := json.Marshal(&modelOrders)
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
	userID, err := getUserID(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := h.user.GetBalance(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("error of get balance")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("current:%v, Withdrawn:%v", current, withdrawn)

	data, err := json.Marshal(&Balance{
		Current:   current,
		Withdrawn: withdrawn,
	})
	if err != nil {
		log.Error().Err(err).Msg("error of serialize balance")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("error of write balance")
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *handler) createWithdraw(rw http.ResponseWriter, r *http.Request) {
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
	log.Printf("createWithdraw, data:%s", string(data))

	var w Withdraw
	err = json.Unmarshal(data, &w)
	if err != nil {
		log.Error().Err(err).Msg("withdraw deserialize error")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("createWithdraw, withdraw:%+v", w)

	err = goluhn.Validate(w.Order)
	if err != nil {
		http.Error(rw, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	orderID, err := strconv.ParseInt(w.Order, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("order number parsing error")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.withdraw.Create(r.Context(), userID, orderID, w.Sum)
	if err != nil {
		if errors.Is(err, withdraw.ErrNotEnoughFunds) {
			rw.WriteHeader(http.StatusPaymentRequired)
			return
		}

		log.Error().Err(err).Msg("error of create withdraw")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

//nolint:dupl //similar to getOrders
func (h *handler) getWithdrawals(rw http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.withdraw.List(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("error of get withdrawals")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("withdrawals:%+v", withdrawals)

	if len(withdrawals) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	modelWithdrawals := make([]WithdrawOut, len(withdrawals))
	for i := 0; i < len(withdrawals); i++ {
		modelWithdrawals[i] = WithdrawOut(withdrawals[i])
	}
	log.Printf("modelWithdrawals:%+v", modelWithdrawals)

	data, err := json.Marshal(&modelWithdrawals)
	if err != nil {
		log.Error().Err(err).Msg("error of serialize withdrawals")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("error of write withdrawals")
		return
	}

	rw.WriteHeader(http.StatusOK)
}
