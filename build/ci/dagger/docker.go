package main

import (
	"context"

	"github.com/lerenn/cryptellation/build/ci/dagger/internal/dagger"
)

type dockerPublisher interface {
	PublishDockerImage(context.Context, *dagger.Directory, []string) error
}

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
	publishers := []dockerPublisher{
		dag.CryptellationBacktestsCi(),
		dag.CryptellationCandlesticksCi(),
		dag.CryptellationExchangesCi(),
		dag.CryptellationForwardtestsCi(),
		dag.CryptellationIndicatorsCi(),
		dag.CryptellationTicksCi(),
	}
	for _, pub := range publishers {
		if err := pub.PublishDockerImage(ctx, sourceDir, tags); err != nil {
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
