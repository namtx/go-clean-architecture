package models

import "errors"

var (
	ErrBadParamInput = errors.New("ErrBadParamInput")
	ErrNotFound      = errors.New("ErrNotFound")
)
