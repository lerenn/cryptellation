package main

import (
	"context"
	"dagger/cryptellation-ci/internal/dagger"
	"fmt"
	"regexp"
	"strings"
)

func updateHelmChartIfNecessary(
	ctx context.Context,
	sourceDir *dagger.Directory,
	repo *Git,
) (*dagger.Directory, error) {
	// Stop here if this not main branch
	// if name, err := repo.GetActualBranch(ctx); err != nil {
	// 	return err
	// } else if name != "main" {
	// 	return nil
	// }

	// Update Helm chart version
	sourceDir, err := updateHelmChartVersion(ctx, sourceDir, repo)
	if err != nil {
		return sourceDir, err
	}

	// Update Helm chart app version
	return updateHelmChartAppVersion(ctx, sourceDir, repo)
}

func updateHelmChartVersion(
	ctx context.Context,
	sourceDir *dagger.Directory,
	repo *Git,
) (*dagger.Directory, error) {
	// Get Helm chart
	helmChart := sourceDir.File("deployments/helm/cryptellation/Chart.yaml")

	// Get content of Helm chart
	content, err := helmChart.Contents(ctx)
	if err != nil {
		return sourceDir, err
	}

	// Compile regexp
	versionRegex, err := regexp.Compile("\nversion: .*")
	if err != nil {
		return sourceDir, err
	}

	// Get version from Helm chart
	version := versionRegex.FindString(content)
	if err != nil {
		return sourceDir, err
	}
	if version == "" {
		return sourceDir, fmt.Errorf("field 'version' not found in Helm chart")
	}
	version = strings.TrimPrefix(version, "\nversion: ")
	version = strings.Trim(version, "\"")

	// Get last commit title
	title, err := repo.GetLastCommitTitle(ctx)
	if err != nil {
		return sourceDir, err
	}

	// Update semver
	newVersion, err := parseAndUpdateSemVer(version, title)
	if err != nil {
		return sourceDir, err
	} else if newVersion == version {
		return sourceDir, nil
	}

	// Update Helm chart
	cmd := "sed -i \"s/^version\\: .*/version\\: " + newVersion + "/\" src/deployments/helm/cryptellation/Chart.yaml"
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

func updateHelmChartAppVersion(
	ctx context.Context,
	sourceDir *dagger.Directory,
	repo *Git,
) (*dagger.Directory, error) {
	// Update app semver
	newVersion, err := repo.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return sourceDir, err
	} else if newVersion == "" {
		return sourceDir, nil
	}

	// Update Helm chart
	cmd := "sed -i \"s/^appVersion\\: .*/appVersion\\: " + newVersion + "/\" src/deployments/helm/cryptellation/Chart.yaml"
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
