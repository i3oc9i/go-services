// Package testgrp maintains the group of handlers for APItesting.
package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/i3oc9i/go-service/business/web/sys/validate"
	"github.com/i3oc9i/go-service/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler is for development.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: "Ok",
	}

	return web.RespondJSON(ctx, w, data, http.StatusOK)
}
