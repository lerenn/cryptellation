package main

import (
	"fmt"
	"strconv"
	"strings"
)

type moduleVersionBump struct {
	Module     string
	Update     string
	OldVersion string
}

func newModuleVersionBump(schema string) moduleVersionBump {
	parts := strings.Split(schema, ":")

	// Avoid errors now
	if len(parts) == 1 {
		parts = append(parts, "")
	}

	return moduleVersionBump{
		Module: parts[0],
		Update: parts[1],
	}
}

func (mvb moduleVersionBump) Validate() error {
	validModule := false

	if mvb.Module == "" {
		validModule = true
	} else {
		for _, name := range modules {
			if name == mvb.Module {
				validModule = true
				break
			}
		}
	}

	if !validModule {
		return fmt.Errorf("unknown module %q, should be one of the following: %v", mvb.Module, modules)
	}

	if mvb.Update != "major" && mvb.Update != "minor" && mvb.Update != "fix" {
		return fmt.Errorf("invalid semver update %q, should be 'major', 'minor' or 'fix'", mvb.Update)
	}

	return nil
}

func (mvb moduleVersionBump) NewVersion() string {
	parts := strings.Split(mvb.OldVersion[1:], ".")

	switch mvb.Update {
	case "major":
		parts[0] = updateStrPart(parts[0])
	case "minor":
		parts[1] = updateStrPart(parts[1])
	case "fix":
		parts[2] = updateStrPart(parts[2])
	default:
		panic("invalid update")
	}

	return "v" + strings.Join(parts, ".")
}

func updateStrPart(part string) string {
	nb, _ := strconv.Atoi(part)
	nb++
	part = strconv.Itoa(nb)
	return part
}

func highestLevelBump(list []moduleVersionBump) string {
	max := "fix"
	for _, mv := range list {
		if mv.Update == "major" {
			max = "major"
		} else if mv.Update == "minor" && max != "major" {
			max = "minor"
		}
	}
	return max
}
