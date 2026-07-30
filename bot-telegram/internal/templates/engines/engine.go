package engine

import (
	"io"
)

type Engine interface {
	Name() string

	Render(wr io.Writer, name string, data interface{}) error

	ParseValues(out *map[string]string, name string, dataString string) error

	Extenstions() []string
}
