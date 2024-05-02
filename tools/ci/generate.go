package main

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/spf13/cobra"
)

func generators() map[string]func(context.Context) error {
	m := make(map[string]func(context.Context) error)
	for _, path := range pathServices {
		m[path] = ci.Generator(client, path)
	}
	return m
}

func runGenerators(cmd *cobra.Command, args []string) {
	ci.ExecuteInParallel(context.Background(), filterWithPath(generators())...)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Execute generator step of the CI",
	Run:     runGenerators,
}

func addGenerateCmdTo(cmd *cobra.Command) {
	cmd.AddCommand(generateCmd)
}
