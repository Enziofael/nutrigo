package v1

import "errors"

var TgIdRequired = errors.New("tg_id is required")
var TgTagRequired = errors.New("tg_tag is required")