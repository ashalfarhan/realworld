package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/go-playground/validator/v10"
)

func errorJSON(w http.ResponseWriter, statusCode int, err error) {
	if statusCode == http.StatusUnprocessableEntity {
		EntityError(w, err)
		return
	}
	JSON(w, statusCode, M{
		"error": err.Error(),
	})
}

func Err(w http.ResponseWriter, e *model.ConduitError) {
	errorJSON(w, e.Code, e.Err)
}

func ClientError(w http.ResponseWriter, err error) {
	errorJSON(w, http.StatusBadRequest, err)
}

func InternalError(w http.ResponseWriter) {
	errorJSON(w, http.StatusInternalServerError, conduit.ErrInternal)
}

func UnauthorizeError(w http.ResponseWriter, reason string) {
	errorJSON(w, http.StatusUnauthorized, fmt.Errorf("%w: %s", conduit.ErrUnauthorized, reason))
}

func EntityError(w http.ResponseWriter, err error) {
	e, ok := err.(validator.ValidationErrors)
	if !ok {
		InternalError(w)
		return
	}

	errors := map[string][]string{}
	for _, field := range e {
		key := strings.ToLower(field.Field())
		errors[key] = append(errors[key], fmt.Sprintf("%s %s", field.Tag(), field.Param()))
	}

	JSON(w, http.StatusUnprocessableEntity, M{
		"errors": errors,
	})
}
