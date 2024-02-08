package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
)

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			telemetry.L(context.TODO()).Error(err.Error())
			panic(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
