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
	jsonBlob, err := json.Marshal(payload)
	if err != nil {
		l.Printf("Error marshalling JSON: %v", err)
		http.Error(w, errMarshallingJSON.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonBlob)
}
