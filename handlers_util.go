package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var (
	errNoAuthHeader    = errors.New("request has no 'Authorization' header")
	errNoBearerToken   = errors.New("'Authorization' header has no bearer token")
	errBadBearerFormat = errors.New("bearer token format is invalid")
	errBadBearerToken  = errors.New("unexpected bearer token")
)

func hasAuth0BearerToken(r *http.Request, expectedToken string) error {
	const prefix = "Bearer "

	value := r.Header.Get("Authorization")
	if value == "" {
		return errNoAuthHeader
	}
	if !strings.HasPrefix(value, prefix) {
		return errNoBearerToken
	}
	elems := strings.Split(value, prefix)
	if len(elems) != 2 {
		return errBadBearerFormat
	}
	if elems[1] != expectedToken {
		return errBadBearerToken
	}

	// All checks have passed and we are confident that the request came from
	// auth0's infrastructure.
	return nil
}

func logAndReturn(w http.ResponseWriter, code int, clientError, logError error) {
	// In case of an error, the extension expects a JSON object containing an
	// "error" attribute.
	jsonBlob, err := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: clientError.Error(),
	})
	if err != nil {
		http.Error(w, `{"error": "failed to marshal error string"}`, http.StatusInternalServerError)
		return
	}

	http.Error(w, string(jsonBlob), code)
	if logError == nil {
		l.Print(clientError)
	} else {
		l.Printf("%v: %v", clientError, logError)
	}
}

func marshalBodyInto(reqBody io.ReadCloser, v interface{}) error {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
