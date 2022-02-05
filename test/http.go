package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

func MakeRequest(method, target string, payload interface{}, handler http.HandlerFunc, bodyDest interface{}) *http.Response {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Panicf("Cannot marshal payload: %v, Reason: %v\n", payload, err)
		return nil
	}

	req := httptest.NewRequest(method, target, bytes.NewReader(b))
	w := httptest.NewRecorder()

	handler(w, req)
	res := w.Result()

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(bodyDest); err != nil {
		log.Panicf("Cannot Decode response body, Reason: %v\n", err)
		return nil
	}

	return res
}
