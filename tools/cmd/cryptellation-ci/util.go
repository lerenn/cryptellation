package main

import (
	"github.com/lerenn/cryptellation/pkg/utils"
)

func filterWithPath[T any](containers map[string]T) []T {
	if pathFlag == "" {
		return utils.MapToList(containers)
	}

	extract := containers[pathFlag]
	return []T{extract}
}
