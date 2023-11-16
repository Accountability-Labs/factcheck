package main

import (
	"encoding/json"
	"io"
	"net/http"
)

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
