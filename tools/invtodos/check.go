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
	err := filepath.Walk(path, func(path string, _ os.FileInfo, err error) error {
		// Return if there is an error
		if err != nil {
			return err
		}

		switch {
		case strings.HasPrefix(path, "tools/invtodos"):
			// Return if this is the path to the tool
			return nil
		case strings.Contains(path, "dagger/internal"):
			// Return if this is a dagger internal file
			return nil
		case filepath.Ext(path) != ".go" &&
			filepath.Ext(path) != ".yaml":
			// Return if it is not: a go file, a yaml file
			return nil
		case strings.HasSuffix(path, ".gen.go"):
			// Return if it is a generated file
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

		// Check if line contains commented TODO
		if !strings.Contains(line, "TODO") {
			continue
		} else if ok := r.MatchString(line); ok {
			continue
		} else if !strings.Contains(line, "//") &&
			!strings.Contains(line, "/*") &&
			!strings.Contains(line, "#") {
			continue
		}

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

	return invalidTodos, nil
}
