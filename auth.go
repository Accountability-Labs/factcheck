package main

import (
	"errors"
	"factcheck/internal/database"
	"net/http"
	"syscall"
)

const (
	authHeader = "X-Auth-Token"
)

var (
	errNoBearerHeader = errors.New("no authentication header")
	errInvalidAPIKey  = errors.New("invalid API key")
	errTalkingToDB    = errors.New("error talking to database")
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
		if errors.Is(err, syscall.ECONNREFUSED) {
			logAndReturn(w, http.StatusServiceUnavailable, errTalkingToDB, err)
			return
		} else if err != nil {
			logAndReturn(w, http.StatusUnauthorized, errInvalidAPIKey, err)
			return
		}

		handler(w, r, &user)
	}
}
