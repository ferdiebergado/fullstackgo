package api

import (
	"encoding/json"
	"io"
)

const contentType = "application/json"

func DecodeJSON(r io.Reader, dest any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return err
	}
	return nil
}
