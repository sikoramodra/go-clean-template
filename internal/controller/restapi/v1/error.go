package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/response"
)

// ErrorResponse -.
func ErrorResponse(w http.ResponseWriter, code int, msg string) {
	WriteJSON(w, code, response.Error{Error: msg})
}
