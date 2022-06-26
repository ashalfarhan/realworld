package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ashalfarhan/realworld/utils/logger"
)

type M map[string]interface{}

func JSON(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log.Printf("Failed to encode json response of %v, Error: %v\n", resp, err)
		fmt.Fprint(w, err.Error())
	}
}
