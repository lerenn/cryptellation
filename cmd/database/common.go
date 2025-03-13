package main

import (
	"github.com/spf13/cobra"
)

func callPersistentPreRunE(cmd *cobra.Command, args []string) error {
	if parent := cmd.Parent(); parent != nil {
		if parent.Context() == nil {
			parent.SetContext(cmd.Context())
		}

		if parent.PersistentPreRunE != nil {
			return parent.PersistentPreRunE(parent, args)
		}
	}

	return nil
}
