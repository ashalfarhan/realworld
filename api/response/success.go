package response

import "net/http"

func Success(w http.ResponseWriter, statusCode int, resp interface{}) {
	JSON(w, statusCode, resp)
}

func Ok(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

func Accepted(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusAccepted, data)
}
