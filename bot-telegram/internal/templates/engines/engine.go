package engine

import (
	"io"
)

type Engine interface {
	Name() string

	Render(wr io.Writer, name string, data interface{}) error

	Extenstions() []string
}
