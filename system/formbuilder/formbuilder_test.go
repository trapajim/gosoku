package formbuilder

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	A string `form:"input,type=text,name=a"`
}
type TestStructNoTags struct {
	A string
}

type TestStructWithTextArea struct {
	A string `form:"input,type=text,name=a"`
	B string `form:"textarea,name=a"`
}
type TestStructInvalidFormType struct {
	A string `form:"invalid,type=text,name=a"`
}

func TestBuild(t *testing.T) {
	type args struct {
		model interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"check form with input text", args{model: TestStruct{}}, `<form> <input type="text" name="a" value=""/> </form>`, false},
		{"check form with input text and value", args{model: TestStruct{A: "hello"}}, `<form> <input type="text" name="a" value="hello"/> </form>`, false},
		{"check form with input text textarea and value value", args{model: TestStructWithTextArea{A: "hello", B: "world"}}, `<form> <input type="text" name="a" value="hello"/>  <textarea name="a">world</textarea> </form>`, false},
		{"check form with input text and textarea", args{model: TestStructWithTextArea{}}, `<form> <input type="text" name="a" value=""/>  <textarea name="a"></textarea> </form>`, false},
		{"pass nil", args{model: nil}, ``, true},
		{"struct without tags", args{model: TestStructNoTags{}}, ``, false},
		{"invalid form type", args{model: TestStructInvalidFormType{}}, ``, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Build(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBuilder_generate(t *testing.T) {
	tests := []struct {
		name    string
		g       DefaultBuilder
		want    string
		wantErr bool
	}{
		{"default builder generate", DefaultBuilder{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := DefaultBuilder{}
			got, err := g.generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultBuilder.generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultBuilder.generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBuilder_getHTML(t *testing.T) {
	tests := []struct {
		name string
		g    DefaultBuilder
		want string
	}{
		{"default builder getHTMl", DefaultBuilder{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := DefaultBuilder{}
			if got := g.getHTML(); got != tt.want {
				t.Errorf("DefaultBuilder.getHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
