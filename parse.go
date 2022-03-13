package tnyrtr

import (
	"encoding/json"
	"net/http"
)

// ParseBodyJSON reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
func ParseBodyJSON(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}
