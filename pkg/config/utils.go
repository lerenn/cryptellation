package config

import (
	"os"
	"strconv"
)

func overrideFromEnv(variable *string, name string) {
	env := os.Getenv(name)
	if name != "" {
		*variable = env
	}
}

func overrideIntFromEnv(variable *int, name string) {
	nb, _ := strconv.Atoi(os.Getenv(name))
	if nb != 0 {
		*variable = nb
	}
}
