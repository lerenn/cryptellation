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
	sourceDir, err := updateHelmChartVersion(ctx, sourceDir, repo, "version")
	if err != nil {
		return sourceDir, err
	}

	// Update Helm chart app version
	return updateHelmChartVersion(ctx, sourceDir, repo, "appVersion")
}

func updateHelmChartVersion(
	ctx context.Context,
	sourceDir *dagger.Directory,
	repo *Git,
	field string,
) (*dagger.Directory, error) {
	// Get Helm chart
	helmChart := sourceDir.File("deployments/helm/cryptellation/Chart.yaml")

	// Get content of Helm chart
	content, err := helmChart.Contents(ctx)
	if err != nil {
		return sourceDir, err
	}

	// Compile regexp
	versionRegex, err := regexp.Compile("\n" + field + ": .*")
	if err != nil {
		return sourceDir, err
	}

	// Get version from Helm chart
	version := versionRegex.FindString(content)
	if err != nil {
		return sourceDir, err
	}
	if version == "" {
		return sourceDir, fmt.Errorf("field %q not found in Helm chart", field)
	}
	version = strings.TrimPrefix(version, "\n"+field+": ")
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
	cmd := "sed -i \"s/^" + field + "\\: .*/" + field + "\\: " + newVersion + "/\" src/deployments/helm/cryptellation/Chart.yaml"
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
