package ui

import (
	"encoding/json"
	"os"
)

// OutputJSON outputs the given data as JSON without wrapping
func OutputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

// OutputError outputs error message as JSON
func OutputError(message string) error {
	errorObj := map[string]string{
		"error": message,
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(errorObj)
}
