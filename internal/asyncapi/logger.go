package asyncapi

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

var _ extensions.Logger = (*LoggerWrapper)(nil)

type LoggerWrapper struct {
}

func (lw LoggerWrapper) Info(ctx context.Context, msg string, info ...extensions.LogInfo) {
	telemetry.L(ctx).Info(mergeInfoWithMessage(msg, info...))
}

func (lw LoggerWrapper) Warning(ctx context.Context, msg string, info ...extensions.LogInfo) {
	telemetry.L(ctx).Warning(mergeInfoWithMessage(msg, info...))
}

func (lw LoggerWrapper) Error(ctx context.Context, msg string, info ...extensions.LogInfo) {
	telemetry.L(ctx).Error(mergeInfoWithMessage(msg, info...))
}

func mergeInfoWithMessage(msg string, info ...extensions.LogInfo) string {
	var mergedInfo string
	for _, in := range info {
		mergedInfo = fmt.Sprintf("%s [%s:%s]", mergedInfo, in.Key, in.Value)
	}

	if len(mergedInfo) > 0 {
		mergedInfo = ":" + mergedInfo
	}

	return fmt.Sprintf("%s%s", msg, mergedInfo)
}
