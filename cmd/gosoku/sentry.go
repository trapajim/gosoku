package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addSentry() error {
	loadDependency()
	dsn := scanSentryDsn()
	addSentryToConfig(dsn)
	createSentryInit()
	addSentryMiddleware()
	return nil
}

func scanSentryDsn() string {
	dsn := ""
	fmt.Print("Sentry dns: ")
	fmt.Scanln(&dsn)
	return dsn
}

func loadDependency() {
	cmd := exec.Command("go", "get", "github.com/getsentry/sentry-go/echo")
	cmd.Run()
}

func addSentryToConfig(dsn string) {
	viper.SetConfigFile("config.yml")
	viper.ReadInConfig() // Find and read the config file
	viper.SetDefault("sentry-dsn", dsn)
	err := viper.WriteConfig()
	fmt.Println(err)
}

func createSentryInit() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(root, "system", "sentry")
	err = os.MkdirAll(path, 0755)
	tmplPath := filepath.Join(root, "cmd", "gosoku", "template", "sentry.tmpl")
	finalFilePath := filepath.Join(path, "sentry.go")
	err = createFile(tmplPath, finalFilePath, domain{})
	return err
}

func addSentryMiddleware() {
	line := FindLine("router.Start()")
	line2 := FindLine("import (")
	if line == 1 {
		fmt.Println("Could not find router.Start()")
		fmt.Println("You should manually append add")

		fmt.Println("router.Echo.Use(sentryecho.New(sentryecho.Options{}))")
		fmt.Println("to main.go")
		return
	}
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectName := getProjectName()
	mainGo := filepath.Join(root, "main.go")
	InsertStringToFile(mainGo, "sentry.InitSentry()\nrouter.Echo.Use(sentryecho.New(sentryecho.Options{}))\n", line-1)
	InsertStringToFile(mainGo, "sentryecho \"github.com/getsentry/sentry-go/echo\"\n\""+projectName+"/system/sentry\"\n", line2)
	cmd := exec.Command("go", "fmt")
	cmd.Run()
}

var sentryCmd = &cobra.Command{
	Use:     "sentry",
	Aliases: []string{"s"},
	Short:   "adds sentry to the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return addSentry()
	},
}

func init() {
	rootCmd.AddCommand(sentryCmd)
}
