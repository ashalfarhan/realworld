package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/go-playground/validator/v10"
)

func ValidateDTO(r *http.Request, dest interface{}) *model.ConduitError {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return conduit.BuildError(400, err)
	}

	v := validator.New()
	if err := v.Struct(dest); err != nil {
		return conduit.BuildError(http.StatusUnprocessableEntity, err)
	}

	return nil
}
