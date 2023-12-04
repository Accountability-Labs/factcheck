package main

import (
	"database/sql"
	"errors"
	"factcheck/internal/database"
	"net/http"
)

var (
	errReadingJSONBody       = errors.New("error reading JSON body")
	errEmailPwdRequired      = errors.New("email and password are required")
	errTalkingToDb           = errors.New("error talking to database")
	errNoteIDAndVoteRequired = errors.New("note ID and vote are required")
	errDerivingPwdHash       = errors.New("error deriving password hash")
	errInvalidEmailPwd       = errors.New("invalid email or password")
)

func (c *apiConfig) getIndex(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hello world."})
}

func (c *apiConfig) getProfile(w http.ResponseWriter, r *http.Request, user *database.User) {
	respondWithJSON(w, http.StatusOK, user)
}

func (c *apiConfig) signinHandler(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}
	if params.Email == "" || params.Password == "" {
		logAndReturn(w, http.StatusBadRequest, errEmailPwdRequired, nil)
		return
	}
	user, err := c.DB.GetUserByEmail(r.Context(), params.Email)
	if errors.Is(err, sql.ErrNoRows) {
		logAndReturn(w, http.StatusUnauthorized, errInvalidEmailPwd, err)
		return
	}
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errTalkingToDb, err)
		return
	}

	hash, err := hashPwdWithSalt(params.Password, user.Salt)
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errDerivingPwdHash, err)
		return
	}
	if hash != user.PasswordHash {
		logAndReturn(w, http.StatusUnauthorized, errInvalidEmailPwd, err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (c *apiConfig) signupHandler(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Email    string `json:"email"`
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}
	if params.Email == "" || params.Password == "" {
		logAndReturn(w, http.StatusBadRequest, errEmailPwdRequired, nil)
		return
	}
	pwdHash, salt, err := hashPwd(params.Password)
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errDerivingPwdHash, err)
		return
	}

	p := database.CreateUserParams{
		Email:        params.Email,
		UserName:     params.UserName,
		PasswordHash: string(pwdHash),
		Salt:         string(salt),
	}
	user, err := c.DB.CreateUser(r.Context(), p)
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

func (c *apiConfig) getRecentNNotes(w http.ResponseWriter, r *http.Request, user *database.User) {
	params := struct {
		Limit int `json:"limit"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}

	notes, err := c.DB.GetRecentNNotes(r.Context(), database.GetRecentNNotesParams{
		VotedBy: user.ID,
		Limit:   5,
	})
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errTalkingToDb, err)
		return
	}
	respondWithJSON(w, http.StatusOK, notes)
}

func (c *apiConfig) getRecentNNotesForUser(w http.ResponseWriter, r *http.Request, user *database.User) {
	params := struct {
		Limit int `json:"limit"`
	}{}

	if err := marshalBodyInto(r.Body, &params); err != nil {
		logAndReturn(w, http.StatusBadRequest, errReadingJSONBody, err)
		return
	}
	notes, err := c.DB.GetRecentNNotesForUser(r.Context(), database.GetRecentNNotesForUserParams{
		ID:    user.ID,
		Limit: int32(params.Limit),
	})
	if err != nil {
		logAndReturn(w, http.StatusInternalServerError, errTalkingToDb, err)
		return
	}
	respondWithJSON(w, http.StatusOK, notes)
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
