package tnyrtr

import (
	"context"
	"encoding/json"
	"net/http"
)

// RespondJSON converts a Go value to JSON and sends it to the client.
func RespondJSON(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	// Set the status code for the request logger middleware.
	err := SetStatusCode(ctx, statusCode)
	if err != nil {
		return err
	}

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
