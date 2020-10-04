package formbuilder

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"
)

const tagName = "form"

// Build creates a new form
func Build(model interface{}) {
	elements := getBuilder(model)
	form := `<form>
	{{ range $key, $value := . }} {{$value.HTML | safeHTML}} {{end}}
	</form>`

	t := template.Must(template.New("base").Funcs(template.FuncMap{
		"safeHTML": func(b string) template.HTML {
			return template.HTML(b)
		},
	}).Parse(form))

	err := t.Execute(os.Stdout, elements)
	fmt.Println(err)

}

// FormElementBuilder is the interface which wraps the basic Generate method.
type FormElementBuilder interface {
	generate() (string, error)
	getHTML() string
}

// DefaultBuilder does not perform any generation.
type DefaultBuilder struct {
}

func (g DefaultBuilder) generate() (string, error) {
	return "", nil
}
func (g DefaultBuilder) getHTML() string {
	return ""
}

// Performs actual generation of the form
func getBuilder(s interface{}) []FormElementBuilder {
	var builders []FormElementBuilder
	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get(tagName)

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		// Get a builder that corresponds to a tag
		builder := getBuilderFromTag(tag)

		// Append error to results
		if builder.getHTML() != "" {
			builders = append(builders, builder)
		}
	}

	return builders
}

func getBuilderFromTag(tag string) FormElementBuilder {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "input":
		var inputType, inputName string
		fmt.Sscanf(strings.Join(args[1:], " "), "type=%s name=%s", &inputType, &inputName)
		builder := NewInput(inputType, inputName)
		return builder
	case "textarea":
		var textareaName string
		fmt.Sscanf(strings.Join(args[1:], " "), "name=%s", &textareaName)
		builder := NewTextarea(textareaName)

		return builder
	}

	return DefaultBuilder{}
}
