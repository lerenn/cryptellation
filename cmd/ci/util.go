package main

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func filterContainerWithPath(containers map[string]*dagger.Container) []*dagger.Container {
	if pathFlag == "" {
		return utils.MapToList(containers)
	}

	extract := containers[pathFlag]
	return []*dagger.Container{extract}
}
