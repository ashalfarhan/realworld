package response

import (
	"net/http"

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

func EntityError(w http.ResponseWriter, err error) {
	e, ok := err.(validator.ValidationErrors)
	if !ok {
		InternalError(w)
		return
	}

	errors := map[string][]string{}
	for _, field := range e {
		errors[field.Field()] = append(errors[field.Field()], field.Error())
	}

	JSON(w, http.StatusUnprocessableEntity, M{
		"errors": errors,
	})
}

func UnauthorizeError(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, conduit.ErrUnauthorized)
}
