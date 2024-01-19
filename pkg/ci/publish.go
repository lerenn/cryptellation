package ci

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/vcs/git"
	"github.com/lerenn/cryptellation/pkg/version"
)

func PublishDockerImage(
	container *dagger.Container,
	serviceName string,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// Publish with hash
		hash, err := version.CommitHashFromGit(".")
		if err != nil {
			return err
		}
		if _, err := container.Publish(ctx, "lerenn/cryptellation-"+serviceName+":"+hash); err != nil {
			return err
		}

		// Stop here if this not main branch
		if name, err := git.ActualBranchName("."); err != nil {
			return err
		} else if name != "main" {
			return nil
		}

		// Publish with version from git
		ver, err := version.VersionFromGit(".", serviceName)
		if err != nil {
			return err
		}
		if _, err := container.Publish(ctx, "lerenn/cryptellation-"+serviceName+":"+ver); err != nil {
			return err
		}

		// Publish with full version
		fullVersion, err := version.FullVersionFromGit(".", serviceName)
		if err != nil {
			return err
		}
		if _, err := container.Publish(ctx, "lerenn/cryptellation-"+serviceName+":"+fullVersion); err != nil {
			return err
		}

		// Publish as latest
		if _, err := container.Publish(ctx, "lerenn/cryptellation-"+serviceName+":latest"); err != nil {
			return err
		}

		return nil
	}
}
