package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/lerenn/cryptellation/pkg/version/git"
	"github.com/spf13/cobra"
)

var (
	modules = []string{
		"cmd/cryptellation",
		"cmd/cryptellation-tui",
		"svc/backtests",
		"svc/candlesticks",
		"svc/exchanges",
		"svc/indicators",
		"svc/ticks",
	}
)

var rootCmd = &cobra.Command{
	Use:   "go run ./tools/tag <path:tag>...",
	Short: "A simple CLI to tag the Cryptellation system",
	Args:  cobra.RangeArgs(1, 100000),
	RunE: func(cmd *cobra.Command, args []string) error {
		mvs, err := toModuleVersion(args...)
		if err != nil {
			return err
		}

		if stop := askMVsConfirmation(mvs...); stop {
			return nil
		}

		return apply(mvs...)
	},
}

func toModuleVersion(args ...string) ([]moduleVersionBump, error) {
	fmt.Print("Generating new module version bump")
	defer fmt.Printf("\n\n")

	// Run over arguments
	list := make([]moduleVersionBump, len(args))
	isGlobal, isAll := false, -1
	for i, arg := range args {
		// Check if this is a '*' argument
		if arg[:2] == "*:" {
			isAll = i
			break
		}

		// Check if this is a codebase update
		if arg[:1] == ":" {
			isGlobal = true
			continue
		}

		// Genrate module version bump
		mvb, err := newModuleVersionBumpFromGit(arg)
		if err != nil {
			return nil, err
		}

		list[i] = mvb
	}

	// If one argument has been detected as '*' for all modules
	if isAll >= 0 {
		list = make([]moduleVersionBump, len(modules))
		for i, m := range modules {
			mv, err := newModuleVersionBumpFromGit(m + args[0][1:])
			if err != nil {
				return nil, err
			}
			list[i] = mv
		}
	}

	// If there is no global version update, generate it
	if !isGlobal { // Get highest bump for the global bump
		max := highestLevelBump(list)

		// Create a global codebase version bump
		globalMVB, err := newModuleVersionBumpFromGit(":" + max)
		if err != nil {
			return nil, err
		}

		list = append(list, globalMVB)
	}

	return list, nil
}

func newModuleVersionBumpFromGit(scheme string) (moduleVersionBump, error) {
	fmt.Print(".")

	// Create a new module version bump from scheme
	mv := newModuleVersionBump(scheme)
	if err := mv.Validate(); err != nil {
		return moduleVersionBump{}, err
	}

	// Get last module version from current branch
	ver, err := git.LastModuleVersionFromCurrentBranch(".", mv.Module)
	if err != nil {
		if !errors.Is(err, git.ErrNoVersion) {
			return moduleVersionBump{}, err
		}
		ver = "v0.0.0"
	}
	mv.OldVersion = ver

	return mv, nil
}

func askMVsConfirmation(mvbs ...moduleVersionBump) bool {
	fmt.Println("Following modules will be bumped:")
	for _, mv := range mvbs {
		module := "[global]"
		if mv.Module != "" {
			module = mv.Module
		}
		fmt.Printf(" %s: %s => %s\n", module, mv.OldVersion, mv.NewVersion())
	}
	fmt.Printf("\n")

	return !askForConfirmation("Do you want to apply these new tags?")
}

func apply(mvs ...moduleVersionBump) error {
	fmt.Println("")
	for _, mv := range mvs {
		// Applying version as module
		if mv.Module != "" {
			fmt.Printf("Applying %s to module %s\n", mv.NewVersion(), mv.Module)
			if err := git.Tag(".", mv.Module+"/"+mv.NewVersion()); err != nil {
				return err
			}
			continue
		}

		// Applying version as codebase
		fmt.Printf("Applying %s to codebase\n", mv.NewVersion())
		if err := git.Tag(".", mv.NewVersion()); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
