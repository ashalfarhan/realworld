package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/go-playground/validator/v10"
)

func Error(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, M{
		"error": err.Error(),
	})
}

func ClientError(w http.ResponseWriter, err error) {
	Error(w, http.StatusBadRequest, err)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, conduit.ErrInternal)
}

func UnauthorizeError(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, conduit.ErrUnauthorized)
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
