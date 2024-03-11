package publish

import (
	"errors"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/version/git"
)

func toModuleVersion(modules, tags []string) ([]moduleVersionBump, error) {
	// Run over arguments
	list := make([]moduleVersionBump, len(tags))
	isGlobal, isAll := false, -1
	for i, arg := range tags {
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
		mvb, err := newModuleVersionBumpFromGit(modules, arg)
		if err != nil {
			return nil, err
		}

		list[i] = mvb
	}

	// If one argument has been detected as '*' for all modules
	if isAll >= 0 {
		list = make([]moduleVersionBump, len(modules))
		for i, m := range modules {
			mv, err := newModuleVersionBumpFromGit(modules, m+tags[0][1:])
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
		globalMVB, err := newModuleVersionBumpFromGit(modules, ":"+max)
		if err != nil {
			return nil, err
		}

		list = append(list, globalMVB)
	}

	return list, nil
}

func newModuleVersionBumpFromGit(modules []string, scheme string) (moduleVersionBump, error) {
	fmt.Printf("Generating versioning for %q\n", scheme)

	// Create a new module version bump from scheme
	mv := newModuleVersionBump(scheme)
	if err := mv.Validate(modules); err != nil {
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

func tagAndPush(mvs ...moduleVersionBump) error {
	for _, mv := range mvs {
		var tag string
		if mv.Module != "" { // Set version as module
			tag = mv.Module + "/" + mv.NewVersion()
		} else { // Set version as codebase
			tag = mv.NewVersion()
		}

		// Apply version
		fmt.Printf("Tagging commit with tag %q\n", tag)
		if err := git.TagCommit(".", tag); err != nil {
			return err
		}

		// Push the tag
		fmt.Printf("Pushing tag %q\n", tag)
		if err := git.PushTags(".", tag); err != nil {
			return err
		}
	}

	return nil
}

func GitTagAndPush(modules, tags []string) error {
	// Stop here if this not main branch
	if name, err := git.ActualBranchName("."); err != nil {
		return err
	} else if name != "main" {
		return nil
	}

	// Validate and compute tags
	mvs, err := toModuleVersion(modules, tags)
	if err != nil {
		return err
	}

	// For each module tag, tag commit and push the tag
	return tagAndPush(mvs...)
}
