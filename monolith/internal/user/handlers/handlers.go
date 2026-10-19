package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spearszx/go-wallet/monolith/internal/api"
	httpErrors "github.com/spearszx/go-wallet/monolith/internal/errors"
	"github.com/spearszx/go-wallet/monolith/internal/user/model"
	"github.com/spearszx/go-wallet/monolith/internal/user/service"
)

func Register(validate *validator.Validate, usrSvc service.UserService) api.HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "CreateUser"

		var request model.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return httpErrors.FromError(err)
		}

		if err := validate.Struct(request); err != nil {
			return httpErrors.FromError(err)
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
		defer cancel()

		user, err := usrSvc.Create(ctx, &request)
		if err != nil {
			return httpErrors.FromError(err)
		}

		return json.NewEncoder(w).Encode(map[string]any{"success": true, "data": user})
	}
}
