package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type M map[string]interface{}

func JSON(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode json response of %v, Error: %s\n", resp, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
