package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func removeFolders(args []string) error {
	name := args[0]
	moduleDirName := strings.ToLower(name)
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	appDir := filepath.Join(root, "app")
	model := filepath.Join(appDir, "domain", moduleDirName+".go")
	module := filepath.Join(appDir, moduleDirName)
	route := filepath.Join(root, "system", "router", moduleDirName+".go")
	os.RemoveAll(model)
	os.RemoveAll(module)
	os.RemoveAll(route)
	return nil
}

var cleanCmd = &cobra.Command{
	Use:     "clean <namespace>",
	Aliases: []string{"c"},
	Short:   "removes the given domain type and it's routes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeFolders(args)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
