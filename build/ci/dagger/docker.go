package main

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/docker"
	"github.com/lerenn/cryptellation/v1/build/ci/dagger/internal/dagger"
)

const (
	dockerImageName = "lerenn/cryptellation"
)

func publishDockerImages(
	ctx context.Context,
	sourceDir *dagger.Directory,
	repo *Git,
) error {
	// Get tags
	tags, err := getDockerTags(ctx, repo)
	if err != nil {
		return err
	}

	// Publish docker images
	if err := publishWorkerDockerImage(ctx, sourceDir, tags); err != nil {
		return err
	}

	return nil
}

// Publishes the worker docker image
func publishWorkerDockerImage(
	ctx context.Context,
	sourceDir *dagger.Directory,
	tags []string,
) error {
	// Get images for each platform
	platformVariants := make([]*dagger.Container, 0, len(docker.GoRunnersInfo))
	for targetPlatform := range docker.GoRunnersInfo {
		runner := dag.Cryptellation().Worker(sourceDir, dagger.CryptellationWorkerOpts{
			TargetPlatform: targetPlatform,
		})

		platformVariants = append(platformVariants, runner)
	}

	// Set publication options from images
	publishOpts := dagger.ContainerPublishOpts{
		PlatformVariants: platformVariants,
	}

	// Publish with tags
	for _, tag := range tags {
		addr := fmt.Sprintf("%s:%s", dockerImageName, tag)
		if _, err := dag.Container().Publish(ctx, addr, publishOpts); err != nil {
			return err
		}
	}

	return nil
}

func getDockerTags(ctx context.Context, repo *Git) ([]string, error) {
	tags := make([]string, 0)

	// Generate last short sha
	lastShortSha, err := repo.GetLastCommitShortSHA(ctx)
	if err != nil {
		return nil, err
	}
	tags = append(tags, lastShortSha)

	// Stop here if this not main branch
	if name, err := repo.GetActualBranch(ctx); err != nil {
		return nil, err
	} else if name != "main" {
		return tags, nil
	}

	// Check if there is a new sem ver, if there is none, just stop here
	newSemVer, err := repo.GetLastTag(ctx)
	if err != nil {
		return nil, err
	} else if newSemVer == "" {
		return tags, nil
	}

	tags = append(tags, newSemVer)
	tags = append(tags, "latest")

	return tags, nil
}
