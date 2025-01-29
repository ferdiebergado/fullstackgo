package api

import (
	"encoding/json"
	"net/http"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
)

type AuthHandler interface {
	HandleUserSignUp(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service service.AuthService
}

var _ AuthHandler = (*authHandler)(nil)

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		service: authService,
	}
}

func (h *authHandler) HandleUserSignUp(w http.ResponseWriter, r *http.Request) {
	var params model.UserSignUpParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.service.SignUpUser(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", contentType)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(userJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
