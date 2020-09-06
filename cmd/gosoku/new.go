package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:     "new [flags] <project name>",
	Short:   "creates a project directory of the name supplied as a parameter",
	Example: `$ gosoku new myproject`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := "gosoku"
		if len(args) > 0 {
			projectName = args[0]
		} else {
			return fmt.Errorf("%s", "No project name supplied")
		}
		return createProject(projectName)
	},
}

func createProject(name string) error {
	projectDir := filepath.Join("app", "domain")
	err := os.MkdirAll(projectDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
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
