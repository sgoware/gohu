package model

import "errors"

var (
	ErrNotSupportGrantType               = errors.New("grant type is not supported")
	ErrNotSupportOperation               = errors.New("no support operation")
	ErrInvalidUsernameAndPasswordRequest = errors.New("invalid username, password")
	ErrInvalidTokenRequest               = errors.New("invalid token")
	ErrExpiredToken                      = errors.New("token is expired")
	ErrInvalidAuthorizationRequest       = errors.New("invalid authorization")
)
