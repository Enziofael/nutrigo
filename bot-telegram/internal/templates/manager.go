package templates

import (
	"bytes"
	"fmt"
	"strings"

	engine "github.com/Enziofael/nutrigo/bot-telegram/internal/templates/engines"
)

type Manager struct {
	engines map[string]engine.Engine
}

func NewManager() (*Manager, error) {
	m := &Manager{
		engines: make(map[string]engine.Engine),
	}

	if err := m.registerEngine(engine.NewHTMLEngine()); err != nil {
		return nil, err
	}

	if err := m.loadTemplates(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) registerEngine(e engine.Engine) error {
	m.engines[e.Name()] = e
	return nil
}

func (m *Manager) loadTemplates() error {
	return nil
}

func (m *Manager) Render(engineName, templateName string, data interface{}) (string, error) {
	eng, ok := m.engines[engineName]

	if !ok {
		return "", fmt.Errorf("Engine %s not found", engineName)
	}

	var buf bytes.Buffer
	if err := eng.Render(&buf, templateName, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (m *Manager) ParseValues(engineName, templateName string, dataString string) (map[string]string, error) {
	eng, ok := m.engines[engineName]

	if !ok {
		return map[string]string{}, fmt.Errorf("Engine %s not found", engineName)
	}

	res := make(map[string]string, 0)
	if err := eng.ParseValues(&res, templateName, dataString); err != nil {
		return map[string]string{}, err
	}

	return res, nil
}

func (m *Manager) RenderHTML(name string, data interface{}) (string, error) {
	return m.Render("html", name, data)
}

func (m *Manager) ParseValuesHTML(name string, data string) (map[string]string, error) {
	return m.ParseValues("html", name, data)
}

func (m *Manager) WrapLineHTML(html string, lineNumber int, openTags, closeTags string) string {
	lines := strings.Split(html, "\n")
	if lineNumber < 1 || lineNumber > len(lines) {
		return html
	}
	lines[lineNumber-1] = openTags + lines[lineNumber-1] + closeTags
	return strings.Join(lines, "\n")
}
