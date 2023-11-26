package main

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/scrypt"
)

const (
	saltLen = 16
	pwdLen  = 32
	N       = 32768
	r       = 8
	p       = 1
)

func hashPwdWithSalt(pwd, salt string) (string, error) {
	rawSalt, err := hex.DecodeString(salt)
	if err != nil {
		return "", err
	}

	rawHash, err := scrypt.Key([]byte(pwd), rawSalt, N, r, p, pwdLen)
	return hex.EncodeToString(rawHash), err
}

func hashPwd(pwd string) (string, string, error) {
	var (
		err     error
		rawSalt = make([]byte, saltLen)
	)
	// Each user has a unique, random salt.
	if _, err = rand.Read(rawSalt); err != nil {
		return "", "", err
	}

	rawHash, err := scrypt.Key([]byte(pwd), rawSalt, N, r, p, pwdLen)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(rawHash), hex.EncodeToString(rawSalt), nil
}
