package main

import (
	"context"

	"github.com/lerenn/cryptellation/internal/dagger/internal/dagger"
)

func (mod *CryptellationInternal) AttachBinance(
	ctx context.Context,
	c *dagger.Container,
	secretsFile *dagger.Secret,
) (*dagger.Container, error) {
	c, err := mod.AttachSecretFromEnvFile(ctx, c, secretsFile, "BINANCE_API_KEY")
	if err != nil {
		return nil, err
	}

	return mod.AttachSecretFromEnvFile(ctx, c, secretsFile, "BINANCE_SECRET_KEY")
}
