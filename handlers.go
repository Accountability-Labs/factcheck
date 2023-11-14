package main

import (
	"errors"
	"factcheck/internal/database"
	"net/http"
)

var (
	errReadingJSONBody       = errors.New("error reading JSON body")
	errEmailRequired         = errors.New("email is required")
	errTalkingToDb           = errors.New("error talking to database")
	errNoteIDAndVoteRequired = errors.New("note ID and vote are required")
)

func (c *apiConfig) getIndex(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hello world."})
}

func (c *apiConfig) getProfile(w http.ResponseWriter, r *http.Request, user *database.User) {
	respondWithJSON(w, http.StatusOK, user)
}

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	params := struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}
	if params.Email == "" {
		logAndReturn(w, http.StatusBadRequest, errEmailRequired, nil)
		return
	}

	user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:    params.Email,
		UserName: params.UserName,
	})
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errTalkingToDb, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (c *apiConfig) createNoteHandler(w http.ResponseWriter, r *http.Request, user *database.User) {
	params := struct {
		Note string `json:"note"`
		Url  string `json:"url"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}

	note, err := c.DB.CreateNote(r.Context(), database.CreateNoteParams{
		CreatedBy: user.ID,
		Note:      params.Note,
		Url:       params.Url,
	})
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errTalkingToDb, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, note)
}

func (c *apiConfig) getRecentNNotesForUrl(w http.ResponseWriter, r *http.Request, user *database.User) {
	params := struct {
		Url   string `json:"url"`
		Limit int    `json:"limit"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}
	// TODO

	notes, err := c.DB.GetRecentNNotesForUrl(r.Context(), database.GetRecentNNotesForUrlParams{
		VotedBy: user.ID,
		Url:     params.Url,
		Limit:   5,
	})
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errTalkingToDb, err)
		return
	}
	respondWithJSON(w, http.StatusOK, notes)
}

func (c *apiConfig) voteOnNote(w http.ResponseWriter, r *http.Request, user *database.User) {
	params := struct {
		NoteID int32 `json:"note_id"`
		Vote   int32 `json:"vote"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}
	if params.NoteID == 0 || params.Vote == 0 {
		logAndReturn(w, http.StatusBadRequest, errNoteIDAndVoteRequired, nil)
		return
	}

	vote, err := c.DB.VoteOnNote(r.Context(), database.VoteOnNoteParams{
		VotedBy: user.ID,
		VotedOn: params.NoteID,
		Vote:    params.Vote,
	})
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errTalkingToDb, err)
		return
	}
	respondWithJSON(w, http.StatusOK, vote)
}
