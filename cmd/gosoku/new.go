package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:     "new [flags] <project name>",
	Short:   "creates a project directory of the name supplied as a parameter",
	Example: `$ gosoku new myproject`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := ""
		fmt.Print("Enter a Projectname (default gosoku):")
		fmt.Scanln(&projectName)
		if projectName == "" {
			projectName = "gosoku"
		}
		return createProject(projectName)
	},
}

func createProject(projectName string) error {
	renameModule(projectName)
	err := createDirectory()
	if err != nil {
		return err
	}
	err = createErrorStruct()
	if err != nil {
		return err
	}
	createConfigYml(projectName)
	return nil
}

func renameModule(moduleName string) {
	cmd := exec.Command("go", "mod", "edit", "-module", moduleName)
	cmd.Run()
}

func createDirectory() error {
	projectDir := filepath.Join("app", "domain")
	err := os.MkdirAll(projectDir, os.ModeDir|os.ModePerm)
	return err
}

func createConfigYml(projectName string) {
	viper.SetConfigName("config")
	viper.SetDefault("name", projectName)
	err := viper.WriteConfig()
	fmt.Println(err)
}

func createErrorStruct() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	tmplPath := filepath.Join(root, "cmd", "gosoku", "template", "err.tmpl")
	componentPathPieces := strings.Split(tmplPath, "/")
	finalFilePath := filepath.Join(root, "app", "domain", "err.go")
	componentFileName := componentPathPieces[len(componentPathPieces)-1]
	tmpl, err := template.New(componentFileName).ParseFiles(tmplPath)
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, nil)
	if err != nil {
		return fmt.Errorf("Failed to execute template: %s", err.Error())
	}
	fmtBuf, err := format.Source(buf.Bytes())
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

func init() {
	rootCmd.AddCommand(newCmd)
}

func execAndWait(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		return err

	}
	return cmd.Wait()
}
