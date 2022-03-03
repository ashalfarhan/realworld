package controller_test

import (
	"net/http"
	"testing"

	. "github.com/ashalfarhan/realworld/api/controller"
	"github.com/ashalfarhan/realworld/test"
	"github.com/stretchr/testify/assert"
)

func TestHelloController(t *testing.T) {
	var result string
	res := test.MakeRequest(http.MethodGet, "/api/hello", nil, &result, Hello)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Contains(t, result, "Hello")
}
