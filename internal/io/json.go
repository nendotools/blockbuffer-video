package io

import (
	"encoding/json"
	"net/http"
)

// SuccessJSON sends a JSON response with a 200 status code
// If only one argument is provided, it assumes a string is a message or else it sends the data
// If two arguments are provided, the first is assumed to be a message and the second is data
func SuccessJSON(w http.ResponseWriter, data ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	jsonData := json.NewEncoder(w)

	if len(data) == 0 {
		jsonData.Encode(map[string]interface{}{"message": "success"})
	}
	if len(data) > 1 {
		jsonData.Encode(map[string]interface{}{"message": data[0], "data": data[1]})
	} else {
		// if data[0] is a string, it's a message
		if _, ok := data[0].(string); ok {
			jsonData.Encode(map[string]interface{}{"message": data[0]})
		} else {
			jsonData.Encode(data[0])
		}
	}
}

func ErrorJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonData := json.NewEncoder(w)
	jsonData.Encode(map[string]interface{}{"error": message, "code": code})
}
