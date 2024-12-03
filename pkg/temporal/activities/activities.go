package activities

import (
	temporalclient "go.temporal.io/sdk/client"
)

// Activities represents common activities that can be used by many workflows.
type Activities struct {
	temporal temporalclient.Client
}

// NewActivities will create new common activities.
func NewActivities(temporal temporalclient.Client) *Activities {
	return &Activities{
		temporal: temporal,
	}
}
