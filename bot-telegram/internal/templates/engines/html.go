package engine

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"strings"
)

//go:embed html_templates/*
var htmlTemplates embed.FS

type HTMLEngine struct {
	templates map[string]*template.Template
}

func NewHTMLEngine() *HTMLEngine {
	e := &HTMLEngine{
		templates: make(map[string]*template.Template),
	}

	err := fs.WalkDir(htmlTemplates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".html") && !strings.HasSuffix(path, ".htm") {
			return nil
		}

		content, err := htmlTemplates.ReadFile(path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(path).Parse(string(content))
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(d.Name(), ".html")
		name = strings.TrimSuffix(name, ".htm")
		e.templates[name] = tmpl

		return nil
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Loaded %d HTML templates:", len(e.templates))
	for name := range e.templates {
		log.Printf("  - %s", name)
	}
	return e
}

func (e *HTMLEngine) Name() string {
	return "html"
}

func (e *HTMLEngine) Extenstions() []string {
	return []string{".html", ".htm"}
}

func (e *HTMLEngine) Render(wr io.Writer, name string, data interface{}) error {
	tmpl, ok := e.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}
	return tmpl.Execute(wr, data)
}