package prompts

import (
	"bytes"
	"fmt"
	"text/template"
)

// Manager handles the loading, rendering, and versioning of AI prompts
type Manager interface {
	// Render executes a named template with the given data
	Render(templateName string, data interface{}) (string, error)
	// RenderSystemPrompt retrieves a system prompt for a specific domain
	RenderSystemPrompt(domain string) (string, error)
}

type manager struct {
	templates      *template.Template
	systemPrompts  map[string]string
}

// NewManager creates a prompt manager loaded with predefined templates
func NewManager() (Manager, error) {
	tmpl := template.New("ai_prompts")
	
	// Load application templates
	for name, content := range applicationTemplates {
		var err error
		tmpl, err = tmpl.New(name).Parse(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
		}
	}

	return &manager{
		templates:     tmpl,
		systemPrompts: systemTemplates,
	}, nil
}

func (m *manager) Render(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := m.templates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}
	return buf.String(), nil
}

func (m *manager) RenderSystemPrompt(domain string) (string, error) {
	prompt, ok := m.systemPrompts[domain]
	if !ok {
		return "", fmt.Errorf("system prompt for domain %s not found", domain)
	}
	return prompt, nil
}
