package main

import (
	"context"
	"dagger/cryptellation-ci/internal/dagger"
	"fmt"
	"regexp"
	"strings"
)

const (
	helmPkgGitRepo = "git@github.com:lerenn/packages.git"
)

func publishHelmChart(
	ctx context.Context,
	sourceDir *dagger.Directory,
	sshPrivateKeyFile *dagger.Secret,
	repo *Git,
) error {
	// Stop here if this not main branch
	if name, err := repo.GetActualBranch(ctx); err != nil {
		return err
	} else if name != "main" {
		return nil
	}

	// Set helm container
	container := dag.Container().From("alpine/helm").
		WithoutEntrypoint().
		WithMountedDirectory("/src", sourceDir).
		WithWorkdir("/")
	if sshPrivateKeyFile != nil {
		container = container.WithMountedSecret("/root/.ssh/id_rsa", sshPrivateKeyFile)
	}

	// Generate package
	container, err := container.
		WithExec([]string{"helm", "package", "/src/deployments/helm/cryptellation"}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Get helm package file name
	entries, err := container.Directory("/").Entries(ctx)
	if err != nil {
		return err
	}
	var pkgFileName string
	for _, entry := range entries {
		pkgFileName = entry
		if strings.HasSuffix(pkgFileName, ".tgz") {
			break
		}
	}

	// Install ssh
	container, err = container.
		WithExec([]string{"apk", "add", "openssh"}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Add SSH config for github
	container, err = container.
		WithExec([]string{"sh", "-c", "echo -e 'Host github.com\n\tStrictHostKeyChecking no\n' > /root/.ssh/config"}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Clone the package repository
	container, err = container.WithExec([]string{"git", "clone", helmPkgGitRepo, "packages"}).Sync(ctx)
	if err != nil {
		return err
	}

	// Move the package to the repository
	container, err = container.WithExec([]string{"mv", "/" + pkgFileName, "/packages/helm/cryptellation"}).Sync(ctx)
	if err != nil {
		return err
	}

	// Update the helm repository index
	container, err = container.
		WithWorkdir("/packages").
		WithExec([]string{
			"helm", "repo", "index",
			"--url", "https://lerenn.github.io/packages/helm/cryptellation",
			"--merge", "./helm/cryptellation/index.yaml",
			"./helm/cryptellation",
		}).Sync(ctx)
	if err != nil {
		return err
	}

	// Add all changes
	container, err = container.
		WithExec([]string{"git", "add", "."}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Add infos on author
	container, err = container.
		WithExec([]string{"git", "config", "--global", "user.email", "louis.fradin+cryptellation-ci@gmail.com"}).
		Sync(ctx)
	if err != nil {
		return err
	}
	container, err = container.
		WithExec([]string{"git", "config", "--global", "user.name", "Cryptellation CI"}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Commit changes
	version := strings.TrimPrefix(pkgFileName, "cryptellation-")
	version = strings.TrimSuffix(version, ".tgz")
	container, err = container.
		WithExec([]string{"git", "commit", "-m", "add cryptellation helm chart " + version}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Push new commit
	_, err = container.
		WithExec([]string{"git", "push"}).
		Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}

func updateHelmChart(
	ctx context.Context,
	sourceDir *dagger.Directory,
	repo *Git,
) (*dagger.Directory, error) {
	// Stop here if this not main branch
	if name, err := repo.GetActualBranch(ctx); err != nil {
		return sourceDir, err
	} else if name != "main" {
		return sourceDir, nil
	}

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
	newVersion = strings.TrimPrefix(newVersion, "v")
	newVersion = fmt.Sprintf("\"%s\"", newVersion)

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
