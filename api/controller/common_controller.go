package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "Hello")
}
