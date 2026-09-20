package v1

import (
	"errors"
	"net/http"

	"github.com/sikoramodra/go-clean-template/internal/controller/restapi/middleware"
	_ "github.com/sikoramodra/go-clean-template/internal/controller/restapi/v1/response"
	"github.com/sikoramodra/go-clean-template/internal/entity"
)

// @Summary     Get profile
// @Description Get current user profile
// @ID          profile
// @Tags        user
// @Produce     json
// @Success     200 {object} entity.User
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /user/profile [get]
func (v1 *V1) profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())

	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	user, err := v1.u.GetUser(r.Context(), userID)
	if err != nil {
		v1.l.Error(err, "restapi - v1 - profile")

		if errors.Is(err, entity.ErrUserNotFound) {
			ErrorResponse(w, http.StatusNotFound, "user not found")

			return
		}

		ErrorResponse(w, http.StatusInternalServerError, "internal server error")

		return
	}

	WriteJSON(w, http.StatusOK, user)
}
