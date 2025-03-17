package worker

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/version"
	"go.temporal.io/sdk/workflow"
)

// ServiceInfo returns the service information.
func ServiceInfo(_ workflow.Context, _ api.ServiceInfoParams) (api.ServiceInfoResults, error) {
	return api.ServiceInfoResults{
		Version: version.Version(),
	}, nil
}
