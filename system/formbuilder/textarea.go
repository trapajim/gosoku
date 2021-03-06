package formbuilder

import (
	"bytes"
	"html/template"
)

// TextareaBuilder generates Textareas for a form
type TextareaBuilder struct {
	Name  string
	Value string
	HTML  string
}

// NewTextarea generates a new Textarea for the form
func NewTextarea(name, value string) TextareaBuilder {
	i := TextareaBuilder{Name: name, Value: value}
	html, err := i.generate()
	if err != nil {
		panic(err)
	}
	i.HTML = html
	return i
}

// Generate creates a template for input fields
func (g TextareaBuilder) generate() (string, error) {
	t, err := template.New("input").Parse(`<textarea name="{{.Name}}">{{.Value}}</textarea>`)
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

func (g TextareaBuilder) getHTML() string {
	return g.HTML
}
