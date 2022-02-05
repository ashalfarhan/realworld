package controller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/test"
)

func TestHelloController(t *testing.T) {
	var result string

	res := test.MakeRequest(http.MethodGet, "/api/hello", nil, Hello, &result)

	json.NewDecoder(res.Body).Decode(&result)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected StatusCode to be %d, but got: %v", http.StatusOK, res.StatusCode)
	}

	if result != "Hello" {
		t.Fatalf("expected Response to be \"%s\", but got: %v", "Hello", result)
	}

}
