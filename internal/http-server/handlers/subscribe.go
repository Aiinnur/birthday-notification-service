package handlers

import (
	resp "birthday-notification-service/internal/appresponse"
	"birthday-notification-service/internal/entity"
	"context"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

type SubscribeRepository interface {
	Subscribe(ctx context.Context, subscription entity.Subscription) error
	Unsubscribe(ctx context.Context, subscription entity.Subscription) error
}

func NewSubscribe(subscribRepo SubscribeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscript entity.Subscription
		err := render.DecodeJSON(r.Body, &subscript)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err = subscribRepo.Subscribe(r.Context(), subscript)
		if err.Error() == fmt.Sprintf("The employee %d already has such a subscriber %d", subscript.BirthdayUserID, subscript.SubscriberID) {
			render.JSON(w, r, resp.Error("this subscription is already exist"))

			return
		}
		if err != nil {
			render.JSON(w, r, resp.Error("failed to create subscription"))

			return
		}

		render.JSON(w, r, resp.Ok())
	}
}

func NewUnsubscribe(subscribRepo SubscribeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var unsubscript entity.Subscription
		err := render.DecodeJSON(r.Body, &unsubscript)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err = subscribRepo.Subscribe(r.Context(), unsubscript)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to delete subscription"))

			return
		}

		render.JSON(w, r, resp.Ok())
	}
}
