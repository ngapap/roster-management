package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"roster-management/internal/models"
)

func SendResponse(w http.ResponseWriter, status int, data, message interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := models.APIResult{
		Data:   data,
		Status: status,
	}

	if message != nil {
		if err, ok := message.(error); ok {
			res.Message = err.Error()
		} else {
			res.Message = fmt.Sprintf("%v", message)
		}
	}

	return json.NewEncoder(w).Encode(res)
}
