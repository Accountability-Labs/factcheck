package main

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"golang.org/x/crypto/scrypt"
)

const (
	saltLen = 16
	pwdLen  = 32
	N       = 32768
	r       = 8
	p       = 1
)

func isValidPort(portStr string) bool {
	maybePort, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return false
	}
	return maybePort > 0 && maybePort <= 65535
}

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
