package service

import "errors"

var ErrTgIdRequired = errors.New("tg_id is required")
var ErrTgTagRequired = errors.New("tg_tag is required")
var ErrInvalidStatus = errors.New("status is invalid")
var ErrInvalidIncludeQuery = errors.New("include is invalid. Allowed values:")