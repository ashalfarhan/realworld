package response

import (
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
)

func Error(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, map[string]string{
		"error": err.Error(),
	})
}

func ClientError(w http.ResponseWriter, err error) {
	Error(w, http.StatusBadRequest, err)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, conduit.ErrInternal)
}

func EntityError(w http.ResponseWriter, err interface{}) {
	JSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
		"errors": err,
	})
}

func UnauthorizeError(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, conduit.ErrUnauthorized)
}
