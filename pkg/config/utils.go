package config

import (
	"os"
)

func overrideString(variable *string, value string) {
	if variable == nil {
		return
	}

	if *variable == "" {
		*variable = value
	}
}

func overrideFromEnv(variable *string, name string) {
	if variable == nil {
		return
	}

	env := os.Getenv(name)
	if env != "" {
		*variable = env
	}
}
