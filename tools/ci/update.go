package main

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/spf13/cobra"
)

func updators() map[string]func(context.Context) error {
	m := make(map[string]func(context.Context) error)
	for _, path := range pathModules {
		m[path] = ci.UpdateGoMod(client, path)
	}
	return m
}

func runUpdators(cmd *cobra.Command, args []string) {
	ci.ExecuteInParallel(context.Background(), filterWithPath(updators())...)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Execute updates step of the CI",
	Run:     runUpdators,
}

func addUpdateCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(updateCmd)
}
