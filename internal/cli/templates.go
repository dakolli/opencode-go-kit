package cli

import tea "github.com/charmbracelet/bubbletea"

type Template struct {
	Name string
}

func newTemplate(name string) Template {
	return Template{
		Name: name,
	}
}

// The message wrapper that will hold your slice of templates
type fetchedTemplatesMsg []Template

func fetchTemplatesCmd() tea.Msg {
	// Simulate fetching templates (replace this with your actual disk read or API call later)
	// e.g., templates := ListLocalTemplates()
	templates := []Template{
		{Name: "Go Backend Template"},
		{Name: "React Frontend Template"},
		{Name: "Python Data Science Template"},
	}

	// Always wrap your raw data in your custom message type before returning
	return fetchedTemplatesMsg(templates)
}
