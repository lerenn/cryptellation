package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func checkInvalidTodosOnDir(path string) ([]string, error) {
	invalidTodos := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Return if there is an error
		if err != nil {
			return err
		}

		// Return if this is the path to the tool
		if strings.HasPrefix(path, "tools/invtodos") {
			return nil
		}

		// Return if it is not a go file
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Return if it is a generated file
		if strings.HasSuffix(path, ".gen.go") {
			return nil
		}

		// Check file for invalid todos
		itd, err := checkInvalidTodosOnFile(path)
		if err != nil {
			return err
		}

		// Append invalid todos to the list
		invalidTodos = append(invalidTodos, itd...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return invalidTodos, nil
}

func checkInvalidTodosOnFile(path string) ([]string, error) {
	r := regexp.MustCompile(`TODO\(#[0-9]+\)`)

	invalidTodos := make([]string, 0)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read file line by line
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check if line contains TODO
		var invalid bool
		if strings.Contains(line, "TODO") {
			if ok := r.MatchString(line); !ok {
				invalid = true
			}
		}

		// Check for invalid todos
		if invalid {
			// Get only the comment
			parts := strings.Split(line, "//")
			if len(parts) > 1 {
				line = strings.Join(parts[1:], "//")
				line = strings.TrimSpace(line)
			}

			// Get the file with the line number
			l := path + ":" + strconv.Itoa(lineNumber)

			// Shorten line
			if len(line) > 20 {
				line = line[:20] + "..."
			}

			description := fmt.Sprintf("%-70s %s", l, line)
			invalidTodos = append(invalidTodos, description)
		}
	}

	return invalidTodos, nil
}
