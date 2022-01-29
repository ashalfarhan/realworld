package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	response.Ok(w, "Hello")
}
