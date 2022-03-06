package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/go-playground/validator/v10"
)

func (m *ConduitMiddleware) WithValidator(next http.HandlerFunc, d interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(d); err != nil {
			response.ClientError(w, err)
			return
		}

		v := validator.New()
		if err := v.Struct(d); err != nil {
			response.EntityError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), dto.DtoCtxKey, d)
		next(w, r.WithContext(ctx))
	}
}
