package config

import (
	"os"
	"strconv"
)

func overrideString(variable *string, value string) {
	if variable == nil {
		return
	}

	if *variable == "" {
		*variable = value
	}
}

func overrideInt(variable *int, value int) {
	if variable == nil {
		return
	}

	if *variable == 0 {
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

func overrideIntFromEnv(variable *int, name string) {
	if variable == nil {
		return
	}

	nb, _ := strconv.Atoi(os.Getenv(name))
	if nb != 0 {
		*variable = nb
	}
}
