package templates

import (
	"bytes"
	"fmt"

	engine "github.com/Enziofael/nutrigo/frontend-bot/internal/templates/engines"
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

func (m *Manager) RenderHTML(name string, data interface{}) (string, error) {
	return m.Render("html", name, data)
}
