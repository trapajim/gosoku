package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type domain struct {
	Name        string
	VarName     string
	ProjectName string // have this here for simplicity
	Fields      []field
}

type field struct {
	Name     string
	VarName  string
	TypeName string
	JSONName string
	EditType string
}

var reservedFields = map[string]string{
	"id":         "ID",
	"slug":       "Slug",
	"created_at": "Timestamp",
	"updated_at": "Updated",
}

func jsonName(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}
func getField(s string, domain *domain) (field, error) {
	// name becomes a string
	// name:string
	// name:string:text
	if !strings.Contains(s, ":") {
		s = s + ":string"
	}
	data := strings.Split(s, ":")

	f := field{
		Name:     data[0],
		VarName:  domain.VarName,
		JSONName: jsonName(data[0]),
		TypeName: data[1],
	}
	f.EditType = "input"
	if len(data) > 2 {
		f.EditType = data[2]
	}
	return f, nil
}

// User name:string email:string
func getType(args []string) (domain, error) {
	t := domain{
		Name:        args[0],
		ProjectName: getProjectName(),
	}
	t.VarName = strings.ToLower(t.Name)
	fields := args[1:]
	for _, field := range fields {
		f, err := getField(field, &t)
		if err != nil {
			return domain{}, err
		}

		t.Fields = append(t.Fields, f)
	}

	if ok, foundConflicts := foundReservedFields(t.Fields); !ok {
		for conflictName := range foundConflicts {
			fmt.Println(fmt.Sprintf("reserved field name: %s", conflictName))
		}

		count := len(foundConflicts)
		var c = "word"
		if count > 1 {
			c = "words"
		}
		return domain{}, fmt.Errorf("You area using (%d) reserved %s", count, c)
	}

	return t, nil
}

func foundReservedFields(fields []field) (bool, map[string]bool) {
	foundConflicts := make(map[string]bool)
	for _, f := range fields {
		if _, ok := reservedFields[strings.ToLower(f.Name)]; ok {
			foundConflicts[f.Name] = true
		}
	}
	if len(foundConflicts) > 1 {
		return false, foundConflicts
	}

	return true, foundConflicts
}

func createModuleFiles(moduleDirName, folder, fileName, template string, gt domain) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	directory := filepath.Join(root, "app", moduleDirName)
	path := filepath.Join(directory, folder)
	err = os.MkdirAll(path, 0755)
	tmplPath := filepath.Join(root, "cmd", "gosoku", "template", template)
	finalFilePath := filepath.Join(path, fileName)
	err = createFile(tmplPath, finalFilePath, gt)
	return err
}
func createRoute(fileName, template string, gt domain) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(root, "system", "router")
	err = os.MkdirAll(path, 0755)
	tmplPath := filepath.Join(root, "cmd", "gosoku", "template", template)
	finalFilePath := filepath.Join(path, fileName)
	err = createFile(tmplPath, finalFilePath, gt)
	return err
}
func generateContentType(args []string) error {
	name := args[0]
	moduleDirName := strings.ToLower(name)
	fileName := moduleDirName + ".go"

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	appDir := filepath.Join(root, "app")
	modelsDir := filepath.Join(appDir, "domain")
	modelsFilePath := filepath.Join(modelsDir, fileName)

	moduleDir := filepath.Join(appDir, moduleDirName)
	// check if module exist
	if _, err := os.Stat(moduleDir); !os.IsNotExist(err) {
		return fmt.Errorf("Module '%s' already exists. Please remove '%s' before running this command", moduleDirName, moduleDir)
	}

	// parse type info from args
	gt, err := getType(args)
	if err != nil {
		return fmt.Errorf("Failed to parse type args: %s", err.Error())
	}
	// create the domain
	tmplPath := filepath.Join(root, "cmd", "gosoku", "template", "domain.tmpl")
	err = createFile(tmplPath, modelsFilePath, gt)
	if err != nil {
		return fmt.Errorf("Failed to parse domain template: %s", err.Error())
	}
	// create folder
	moduleDirectoryPath := filepath.Join(root, "app", moduleDirName)
	err = os.MkdirAll(moduleDirectoryPath, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %s", err.Error())
	}
	// create repository
	err = createModuleFiles(moduleDirName, "repository", "repository.go", "repository_psql.tmpl", gt)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %s", err.Error())
	}
	// create usecase
	err = createModuleFiles(moduleDirName, "usecase", "usecase.go", "usecase.tmpl", gt)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %s", err.Error())
	}
	// create delivery json
	err = createModuleFiles(moduleDirName, "delivery", "json.go", "delivery_json.tmpl", gt)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %s", err.Error())
	}

	err = createRoute(gt.VarName+".go", "route_init.tmpl", gt)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %s", err.Error())
	}
	return nil
}

func createFile(tmplPath string, finalFilePath string, gt domain) error {
	funcMap := template.FuncMap{
		"add": func(i int) int {
			return i + 1
		},
	}
	componentPathPieces := strings.Split(tmplPath, "/")

	// Get the last item in the pieces (this should be the file name).
	componentFileName := componentPathPieces[len(componentPathPieces)-1]
	tmpl, err := template.New(componentFileName).Funcs(funcMap).ParseFiles(tmplPath)

	if _, err := os.Stat(finalFilePath); !os.IsNotExist(err) {
		return fmt.Errorf("Please remove '%s' before running this command", finalFilePath)
	}
	if err != nil {
		return fmt.Errorf("Failed to parse template: %s", err.Error())
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, gt)
	if err != nil {
		return fmt.Errorf("Failed to execute template: %s", err.Error())
	}
	fmtBuf, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to format template: %s", err.Error())
	}

	// create model
	file, err := os.Create(finalFilePath)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.Write(fmtBuf)
	if err != nil {
		return fmt.Errorf("Failed to generated file buffer: %s", err.Error())
	}
	return nil
}

var generateCmd = &cobra.Command{
	Use:     "generate <generator type (,...fields)>",
	Aliases: []string{"g"},
	Short:   "generate boilerplate code for an api endpoint",
}

var scaffoldCmd = &cobra.Command{
	Use:     "scaffold <namespace> <field> <field>...",
	Aliases: []string{"s"},
	Short:   "auto-generation of a set of a model, routes, and a controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateContentType(args)
	},
}

func init() {
	generateCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(generateCmd)
}
