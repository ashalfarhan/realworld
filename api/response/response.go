package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err != nil {
		data = map[string]interface{}{
			"errors": err,
		}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode json response of %v, Error: %s\n", data, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	JSON(w, statusCode, data, nil)
}

func Error(w http.ResponseWriter, statusCode int, err interface{}) {
	JSON(w, statusCode, nil, err)
}
