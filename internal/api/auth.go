package api

import (
	"net/http"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
)

type AuthHandler interface {
	HandleUserSignUp(w http.ResponseWriter, r *http.Request)
	HandleUserSignIn(w http.ResponseWriter, r *http.Request)
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
	if err := DecodeJSON(r, &params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.service.SignUpUser(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, http.StatusCreated, newUser)
}

func (h *authHandler) HandleUserSignIn(w http.ResponseWriter, r *http.Request) {
	var params model.UserSignInParams
	if err := DecodeJSON(r, &params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.service.SignInUser(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, http.StatusOK, []byte("{}"))
}
