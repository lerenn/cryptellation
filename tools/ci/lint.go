package main

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/spf13/cobra"
)

func linters() map[string]*dagger.Container {
	m := make(map[string]*dagger.Container)
	for _, path := range pathModules {
		m[path] = ci.Linter(client, path)
	}
	return m
}

func runLinters(cmd *cobra.Command, args []string) {
	ci.ExecuteContainersInParallel(
		context.Background(),
		filterWithPath(linters()),
	)
}

var lintCmd = &cobra.Command{
	Use:     "lint",
	Aliases: []string{"l"},
	Short:   "Execute linter step of the CI",
	Run:     runLinters,
}

func addLintCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(lintCmd)
}
