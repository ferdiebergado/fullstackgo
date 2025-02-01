package handler

import (
	"errors"
	"net/http"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/pkg/validation"
	"github.com/ferdiebergado/fullstackgo/internal/service"
	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	HandleUserSignUp(w http.ResponseWriter, r *http.Request)
	HandleUserSignIn(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service   service.AuthService
	validator validation.Validator
}

var _ AuthHandler = (*authHandler)(nil)

func NewAuthHandler(authService service.AuthService, validator validation.Validator) AuthHandler {
	return &authHandler{
		service:   authService,
		validator: validator,
	}
}

func (h *authHandler) HandleUserSignUp(w http.ResponseWriter, r *http.Request) {
	var params model.UserSignUpParams
	if err := DecodeJSON(r, &params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(params); err != nil {
		var valErrs *validator.ValidationErrors
		if errors.As(err, &valErrs) {
			responseValError(w, valErrs)
			return
		}

		serverError(w)
		return
	}

	user, err := h.service.SignUpUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			res := APIResponse{
				Message: "Invalid input!",
				Errors: []map[string]string{
					{"email": err.Error()},
				},
			}

			responseJSON(w, http.StatusUnprocessableEntity, res)
			return
		}

		serverError(w)
		return
	}

	responseJSON(w, http.StatusCreated, user)
}

func (h *authHandler) HandleUserSignIn(w http.ResponseWriter, r *http.Request) {
	var params model.UserSignInParams
	if err := DecodeJSON(r, &params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(params); err != nil {
		var valErrs *validator.ValidationErrors
		if errors.As(err, &valErrs) {
			responseValError(w, valErrs)
			return
		}

		serverError(w)
		return
	}

	_, err := h.service.SignInUser(r.Context(), params)
	if err != nil {
		serverError(w)
		return
	}

	responseJSON(w, http.StatusOK, []byte("{}"))
}
