package engine

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
)

//go:embed html_templates/*
var htmlTemplates embed.FS

type HTMLEngine struct {
	templates    map[string]*template.Template
	rawTemplates map[string][]byte
}

func NewHTMLEngine() *HTMLEngine {
	e := &HTMLEngine{
		templates:    make(map[string]*template.Template),
		rawTemplates: make(map[string][]byte),
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
		e.rawTemplates[name] = content

		return nil
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Loaded %d HTML templates:", len(e.templates))
	for name := range e.templates {
		log.Printf("  - %s", name)
	}
	fmt.Print("\n")
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
		file, _ := os.Create(fmt.Sprintf("F:/Kuroshio/Proging/repos/nutrigo/bot-telegram/internal/templates/engines/html_templates/_%s.html", name))

		file.Close()
		return fmt.Errorf("template %s not found. ", name)
	}
	return tmpl.Execute(wr, data)
}

func (e *HTMLEngine) ParseValues(out *map[string]string, name string, dataString string) error {
	rawTmpl, ok := e.rawTemplates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}
	return Parse(string(rawTmpl), out, dataString)
}

func Parse(rawTmpl string, out *map[string]string, dataString string) error {

	pattern := regexp.QuoteMeta(rawTmpl)

	rePlaceholder := regexp.MustCompile(`{{\s*\.(\w+)\s*}}`)
	matches := rePlaceholder.FindAllStringSubmatch(pattern, -1)

	if len(matches) == 0 {
		*out = map[string]string{}
		return nil
	}

	for _, m := range matches {
		fullMatch := m[0]
		fieldName := m[1]
		quotedFull := regexp.QuoteMeta(fullMatch)
		pattern = strings.Replace(pattern, quotedFull, `(?P<`+fieldName+`>.+)`, 1)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("compile regex: %w", err)
	}

	match := re.FindStringSubmatch(dataString)
	if match == nil {
		return fmt.Errorf("rendered text does not match template")
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	*out = result
	return nil
}

