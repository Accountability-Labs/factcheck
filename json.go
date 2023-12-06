package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errMarshallingJSON = errors.New("error marshalling json")
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// The JSON response body is always available via the "data" attribute.
	wrapper := struct {
		Data interface{} `json:"data"`
	}{
		Data: payload,
	}
	jsonBlob, err := json.Marshal(wrapper)
	if err != nil {
		l.Printf("Error marshalling JSON: %v", err)
		http.Error(w, errMarshallingJSON.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonBlob)
}
