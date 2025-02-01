package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Message string              `json:"message"`
	Errors  []map[string]string `json:"errors,omitempty"`
	Data    any                 `json:"data,omitempty"`
}

func DecodeJSON(r *http.Request, dest any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

func responseJSON(w http.ResponseWriter, status int, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func responseValError(w http.ResponseWriter, valErrs *validator.ValidationErrors) {
	errs := make([]map[string]string, 0)
	for _, e := range *valErrs {
		errs = append(errs, map[string]string{
			e.Field(): e.Error(),
		})
	}

	res := APIResponse{
		Message: "Invalid input!",
		Errors:  errs,
	}

	responseJSON(w, http.StatusUnprocessableEntity, res)
}

func serverError(w http.ResponseWriter) {
	http.Error(w, "An error occurred.", http.StatusInternalServerError)
}

func handleValidationError(w http.ResponseWriter, err error) bool {
	var valErrs *validator.ValidationErrors
	if errors.As(err, &valErrs) {
		responseValError(w, valErrs)
		return true
	}

	serverError(w)
	return true
}
