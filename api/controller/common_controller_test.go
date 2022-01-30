package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloController(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	w := httptest.NewRecorder()
	Hello(w, req)
	res := w.Result()
	defer res.Body.Close()

	expected := map[string]interface{}{
		"result": "Hello",
		"status": 200,
	}
	var result string
	json.NewDecoder(res.Body).Decode(&result)

	if res.StatusCode != expected["status"] {
		t.Fatalf("expected StatusCode to be %d, but got: %v", expected["status"], result)
	}

	if result != expected["result"] {
		t.Fatalf("expected Response to be \"%s\", but got: %v", expected["result"], result)
	}

}
