package publish

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/version"
	"github.com/lerenn/cryptellation/pkg/version/git"
)

func PublishDockerImage(
	container *dagger.Container,
	moduleName, imageName string,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// Publish with hash
		hash, err := version.CommitHashFromGit(".")
		if err != nil {
			return err
		}
		if _, err := container.Publish(ctx, imageName+":"+hash); err != nil {
			return err
		}

		// Stop here if this not main branch
		if name, err := git.ActualBranchName("."); err != nil {
			return err
		} else if name != "main" {
			return nil
		}

		// Publish with version from git
		ver, err := version.VersionFromGit(".", moduleName)
		if err != nil {
			return err
		}
		if _, err := container.Publish(ctx, imageName+":"+ver); err != nil {
			return err
		}

		// Publish with full version
		fullVersion, err := version.FullVersionFromGit(".", moduleName)
		if err != nil {
			return err
		}
		if _, err := container.Publish(ctx, imageName+":"+fullVersion); err != nil {
			return err
		}

		// Publish as latest
		if _, err := container.Publish(ctx, imageName+":latest"); err != nil {
			return err
		}

		return nil
	}
}
