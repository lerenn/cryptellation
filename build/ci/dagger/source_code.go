package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/lerenn/cryptellation/v1/build/ci/dagger/internal/dagger"
)

func updateSourceCode(ctx context.Context, sourceDir *dagger.Directory, repo *Git) (*dagger.Directory, error) {
	// Stop here if this not main branch
	if name, err := repo.GetActualBranch(ctx); err != nil {
		return sourceDir, err
	} else if name != "main" {
		return sourceDir, nil
	}

	// Get new version and check what type of change it is
	change, _, err := repo.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return sourceDir, err
	} else if change == SemVerChangeNone {
		return sourceDir, nil
	}

	// Update version package

	// Update source code based on specific change
	switch change {
	case SemVerChangeMajor:
		return updateSourceCodeWithMajor(ctx, sourceDir, repo)
	case SemVerChangeMinor:
		fallthrough // No change needed
	case SemVerChangePatch:
		fallthrough // No change needed
	default:
		return sourceDir, nil // No change needed
	}
}

func updateVersionPackage(ctx context.Context, sourceDir *dagger.Directory, repo *Git) (*dagger.Directory, error) {
	// Get new version
	change, ver, err := repo.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return sourceDir, err
	}

	// Return if there is nothing to update
	if change == SemVerChangeNone {
		return sourceDir, nil
	}

	// Update version package
	cmd := fmt.Sprintf("sed -i 's/^	DefaultVersion = .*/	DefaultVersion = %s/g' src/pkg/version/version.go", ver)
	c, err := dag.Container().From("alpine").
		WithMountedDirectory("src", sourceDir).
		WithExec([]string{"sh", "-c", cmd}).
		Sync(ctx)
	if err != nil {
		return sourceDir, err
	}

	// Export modified directory
	return c.Directory("src"), nil
}

func updateSourceCodeWithMajor(ctx context.Context, sourceDir *dagger.Directory, repo *Git) (*dagger.Directory, error) {
	// Get new semver
	change, newVersion, err := repo.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return sourceDir, err
	}

	// Return if there is nothing to update
	if change == SemVerChangeNone {
		return sourceDir, nil
	}

	// Get new major version
	newMajor, err := strconv.Atoi(strings.Split(newVersion, ".")[0])
	if err != nil {
		return sourceDir, fmt.Errorf("could not parse new major version: %w", err)
	}
	oldMajor := newMajor - 1

	// Update Cryptellation import directives
	fileTypeArg := "-type f \\( -iname \\*.go -o -iname \\*.go.mod \\)"
	oldPath := fmt.Sprintf("github.com\\/lerenn\\/cryptellation\\/v%d", oldMajor)
	newPath := fmt.Sprintf("github.com\\/lerenn\\/cryptellation\\/v%d", newMajor)
	cmd := fmt.Sprintf("find . %s -exec sed -i 's/%s/%s/g' {} \\;", fileTypeArg, oldPath, newPath)
	c, err := dag.Container().From("alpine").
		WithMountedDirectory("src", sourceDir).
		WithExec([]string{"sh", "-c", cmd}).
		Sync(ctx)
	if err != nil {
		return sourceDir, err
	}

	// Export modified directory
	return c.Directory("src"), nil
}
