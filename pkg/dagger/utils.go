package main

import (
	"context"
	"cryptellation/pkg/dagger/internal/dagger"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

func (mod *CryptellationPkg) WithGoCodeAndCacheAsWorkDirectory(
	c *dagger.Container,
	rootDir *dagger.Directory,
	path string,
) *dagger.Container {
	return c.
		// Add Go caches
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("gobuild")).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gocache")).

		// Add source code
		WithMountedDirectory("/go/src/cryptellation", rootDir).

		// Add workdir
		WithWorkdir("/go/src/cryptellation/" + path)
}

func (mod *CryptellationPkg) CryptellationGoCodeContainer(
	sourceDir *dagger.Directory,
	path string,
) *dagger.Container {
	c := dag.Container().From(golangImage)
	return mod.WithGoCodeAndCacheAsWorkDirectory(c, sourceDir, path)
}

func (mod *CryptellationPkg) LoadSecretFromEnvFile(
	ctx context.Context,
	secretFile *dagger.Secret,
	name string,
) (*dagger.Secret, error) {
	// Load secret file
	plain, err := secretFile.Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	// Load file with secrets
	envMap, err := godotenv.Parse(strings.NewReader(plain))
	if err != nil {
		return nil, err
	}

	// Get requested secret from loaded file
	content, exists := envMap[name]
	if !exists {
		return nil, fmt.Errorf("secret %s not found in secret file", name)
	}

	// Change to dagger secret
	return dag.SetSecret(name, content), nil
}

func (mod *CryptellationPkg) AttachSecretFromEnvFile(
	ctx context.Context,
	c *dagger.Container,
	secretFile *dagger.Secret,
	name string,
) (*dagger.Container, error) {
	secret, err := mod.LoadSecretFromEnvFile(ctx, secretFile, name)
	if err != nil {
		return nil, err
	}

	return c.WithSecretVariable(name, secret), nil
}
