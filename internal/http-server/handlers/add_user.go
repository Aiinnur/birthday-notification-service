package handlers

import (
	resp "birthday-notification-service/internal/appresponse"
	"birthday-notification-service/internal/entity"
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/render"
	"net/http"
)

type UserRepository interface {
	AddUser(ctx context.Context, user entity.User) error
}

func New(userRepo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user entity.User
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		if !govalidator.IsEmail(user.Email) {
			render.JSON(w, r, resp.Error("invalid email in request"))

			return
		}

		err = userRepo.AddUser(r.Context(), user)
		if err.Error() == fmt.Sprintf("user with email %s already exists", user.Email) {
			render.JSON(w, r, resp.Error("a user with this email already exist"))

			return
		}
		if err != nil {
			render.JSON(w, r, resp.Error("failed to add user"))

			return
		}

		render.JSON(w, r, resp.Ok())
	}
}
