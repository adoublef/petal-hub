package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
)

const (
	contentTypJSON = "application/json"
)

var (
	ErrBodyNotFound       = errors.New("http: request body is missing")
	ErrContentTypNotFound = errors.New("http: content-type header not found")
	ErrContentTypJSON     = errors.New("http: content-type is not application/json")
)

// decode JSON Body of incoming request
func decode[T any](_ http.ResponseWriter, r *http.Request) (T, error) {
	var v T
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		return v, ErrContentTypNotFound
	}
	if r.Body == nil {
		return v, ErrBodyNotFound
	}
	d, _, err := mime.ParseMediaType(ct)
	if err != nil || !(d == contentTypJSON) {
		return v, ErrContentTypJSON
	}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decoding json: %w", err)
	}
	return v, nil
}

// encode to the client using JSON
func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return nil
	}
	w.Header().Set("Content-Type", contentTypJSON)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
