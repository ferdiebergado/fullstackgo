package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const contentType = "application/json"

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

	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
