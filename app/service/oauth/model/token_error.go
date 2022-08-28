package model

import "errors"

var (
	ErrNotSupportGrantType         = errors.New("grant type is not supported")
	ErrNotSupportOperation         = errors.New("no support operation")
	ErrInvalidUserId               = errors.New("invalid user id")
	ErrInvalidToken                = errors.New("invalid token")
	ErrExpiredToken                = errors.New("token is expired")
	ErrInvalidAuthorizationRequest = errors.New("invalid authorization")
	ErrUserDetailNotFound          = errors.New("user details not found")
)
