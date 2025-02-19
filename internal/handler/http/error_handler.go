package http

import (
	"context"
	"github.com/DenisBarabanshchikov/subscription/internal/handler/http/response"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"net/http"
)

func handleError(ctx context.Context, err error) response.ErrorResponse {
	switch e := err.(type) {
	case model.CustomerNotFoundErr, model.SubscriptionNotFoundErr:
		return response.ErrorResponse{Code: http.StatusNotFound, Message: e.Error()}
	default:
		return response.ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}
