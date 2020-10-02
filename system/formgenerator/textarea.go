package formgenerator

import (
	"bytes"
	"html/template"
)

// TextareaGenerator generates Textareas for a form
type TextareaGenerator struct {
	Name string
	HTML string
}

// NewTextarea generates a new Textarea for the form
func NewTextarea(name string) TextareaGenerator {
	i := TextareaGenerator{Name: name}
	html, err := i.generate()
	if err != nil {
		panic(err)
	}
	i.HTML = html
	return i
}

// Generate creates a template for input fields
func (g TextareaGenerator) generate() (string, error) {
	t, err := template.New("input").Parse(`<textarea name="{{.Name}}"></textarea>`)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, g)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g TextareaGenerator) getHTML() string {
	return g.HTML
}
