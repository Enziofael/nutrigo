package repository

import "errors"

var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
var ErrUnmarshalJSON = errors.New("can't unmarshall JSON")
var ErrMarshalJSON = errors.New("can't marshal JSON")
var ErrInvalidRequest = errors.New("invalid request")
var ErrRequestConflict = errors.New("request conflicts with the other one")
