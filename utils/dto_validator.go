package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ashalfarhan/realworld/model"
	"github.com/go-playground/validator/v10"
)

func ValidateDTO(r *http.Request, dest interface{}) *model.ConduitError {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return &model.ConduitError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}

	v := validator.New()
	if err := v.Struct(dest); err != nil {
		return &model.ConduitError{
			Code: http.StatusUnprocessableEntity,
			Err:  err,
		}
	}

	return nil
}
