package main

import (
	"context"
	"dagger/cryptellation-ci/internal/dagger"
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

	if err := dag.CryptellationBacktestsCi().PublishDockerImage(ctx, sourceDir, tags); err != nil {
		return err
	}

	if err := dag.CryptellationCandlesticksCi().PublishDockerImage(ctx, sourceDir, tags); err != nil {
		return err
	}

	if err := dag.CryptellationExchangesCi().PublishDockerImage(ctx, sourceDir, tags); err != nil {
		return err
	}

	if err := dag.CryptellationForwardtestsCi().PublishDockerImage(ctx, sourceDir, tags); err != nil {
		return err
	}

	if err := dag.CryptellationIndicatorsCi().PublishDockerImage(ctx, sourceDir, tags); err != nil {
		return err
	}

	if err := dag.CryptellationTicksCi().PublishDockerImage(ctx, sourceDir, tags); err != nil {
		return err
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
	newSemVer, err := repo.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return nil, err
	} else if newSemVer == "" {
		return tags, nil
	}

	tags = append(tags, newSemVer)
	tags = append(tags, "latest")

	return tags, nil
}
