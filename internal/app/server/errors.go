package server

import "errors"

var (
	errInvalidToken       = errors.New("invalid token")
	errWrongSingingMethod = errors.New("unexpected signing method")
)
