package response

import (
	"encoding/json"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
)

type M map[string]interface{}

func JSON(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		conduit.Logger.Printf("Failed to encode json response of %v, Error: %v\n", resp, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
