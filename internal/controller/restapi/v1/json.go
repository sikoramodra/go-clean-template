package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
)

// DecodeAndValidate -.
func DecodeAndValidate(r *http.Request, v *validator.Validate, data any) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}

	if err := v.Struct(data); err != nil {
		return err
	}

	return nil
}

// WriteJSON -.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	// TODO: handle error
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}
