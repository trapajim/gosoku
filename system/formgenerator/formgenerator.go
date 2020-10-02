package formgenerator

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"
)

const tagName = "form"

// Create creates a new form
func Create(model interface{}) {
	elements := getGenerators(model)
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

// FormElementGenerator is the interface which wraps the basic Generate method.
type FormElementGenerator interface {
	generate() (string, error)
	getHTML() string
}

// DefaultGenerator does not perform any generation.
type DefaultGenerator struct {
}

func (g DefaultGenerator) generate() (string, error) {
	return "", nil
}
func (g DefaultGenerator) getHTML() string {
	return ""
}

// Performs actual data validation using validator definitions on the struct
func getGenerators(s interface{}) []FormElementGenerator {
	var generators []FormElementGenerator
	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get(tagName)

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		// Get a validator that corresponds to a tag
		generator := getValidatorFromTag(tag)

		// Append error to results
		if generator.getHTML() != "" {
			generators = append(generators, generator)
		}
	}

	return generators
}

func getValidatorFromTag(tag string) FormElementGenerator {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "input":
		var inputType, inputName string
		fmt.Sscanf(strings.Join(args[1:], " "), "type=%s name=%s", &inputType, &inputName)
		generator := NewInput(inputType, inputName)
		return generator
	case "textarea":
		var textareaName string
		fmt.Sscanf(strings.Join(args[1:], " "), "name=%s", &textareaName)
		generator := NewTextarea(textareaName)

		return generator
	}

	return DefaultGenerator{}
}
