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

func removeLeadingSlash(list []string) []string {
	for i, s := range list {
		if len(s) > 0 && s[0] == '/' {
			list[i] = s[1:]
		}
	}
	return list
}
