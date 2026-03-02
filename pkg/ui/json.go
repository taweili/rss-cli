package ui

import (
	"encoding/json"
	"os"
)

// OutputJSON outputs the given data as JSON without wrapping
// Deprecated: Use Printer.Output() instead for consistent format support
func OutputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

// OutputError outputs error message as JSON
// Deprecated: Use Printer.Error() instead for consistent format support
func OutputError(message string) error {
	errorObj := map[string]string{
		"error": message,
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(errorObj)
}
