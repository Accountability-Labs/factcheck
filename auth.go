package main

import (
	"errors"
	"factcheck/internal/database"
	"net/http"
)

const (
	authHeader = "X-Auth-Token"
)

var (
	errNoBearerHeader = errors.New("no authentication header")
	errInvalidAPIKey  = errors.New("invalid API key")
)

type authHandler func(w http.ResponseWriter, r *http.Request, user *database.User)

func (c *apiConfig) authenticate(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get(authHeader)
		if apiKey == "" {
			http.Error(w, errNoBearerHeader.Error(), http.StatusUnauthorized)
			return
		}

		user, err := c.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			l.Printf("Error retrieving user by API key: %v", err)
			http.Error(w, errInvalidAPIKey.Error(), http.StatusUnauthorized)
			return
		}

		handler(w, r, &user)
	}
}
