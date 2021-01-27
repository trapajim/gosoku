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

type dBData struct {
	Port     int
	Host     string
	User     string
	Password string
	DbName   string
}

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

func scanDatabaseData(projectName string) dBData {
	db := dBData{}
	fmt.Println(projectName)
	fmt.Println(string("\033[36m"), "******************")
	fmt.Println("", "*DB Configuration*")
	fmt.Println("", "******************", string("\033[0m"))
	fmt.Print("Port (default 5432):")
	fmt.Scanln(&db.Port)
	if db.Port == 0 {
		db.Port = 5432
	}
	fmt.Print("Host (default localhost):")
	fmt.Scanln(&db.Host)
	if db.Host == "" {
		db.Host = "localhost"
	}
	fmt.Printf("Database name (default %s):", projectName)
	fmt.Scanln(&db.DbName)
	if db.DbName == "" {
		db.DbName = projectName
	}
	fmt.Printf("User (default %s):", projectName)
	fmt.Scanln(&db.User)
	if db.User == "" {
		db.User = projectName
	}
	scanPassword(&db)
	return db
}

func scanPassword(db *dBData) {
	fmt.Print("Password:")
	fmt.Scanln(&db.Password)
	if db.Password == "" {
		fmt.Println(string("\033[31m"), "Password can't be empty", string("\033[0m"))
		scanPassword(db)
	}
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
	data := scanDatabaseData(projectName)
	fmt.Println(data)
	createConfigYml(projectName, data)
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

func createConfigYml(projectName string, db dBData) {
	viper.SetConfigName("config")
	viper.SetDefault("name", projectName)
	viper.SetDefault("database", db)
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
