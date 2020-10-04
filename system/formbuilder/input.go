package formbuilder

import (
	"bytes"
	"html/template"
)

// InputBuilder generates input fields for the a form
type InputBuilder struct {
	Type string
	Name string
	HTML string
}

// NewInput generates a new Input element for the form
func NewInput(inputType, name string) InputBuilder {
	i := InputBuilder{Name: name, Type: inputType}
	html, err := i.generate()
	if err != nil {
		panic(err)
	}
	i.HTML = html
	return i
}

// Generate creates a template for input fields
func (g InputBuilder) generate() (string, error) {
	t, err := template.New("input").Parse(`<input type="{{.Type}}" name="{{.Name}}" />`)
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

func (g InputBuilder) getHTML() string {
	return g.HTML
}
