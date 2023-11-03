package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func logAndReturn(w http.ResponseWriter, code int, clientError, logError error) {
	http.Error(w, clientError.Error(), code)
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
