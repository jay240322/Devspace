package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var serviceNamePattern = regexp.MustCompile(`^[a-z]([a-z0-9-]*[a-z0-9])?$`)
var reservedNamePattern = regexp.MustCompile(`^(con|prn|aux|nul|com[1-9]|lpt[1-9])$`)

func resolveProjectDirectory(input string) (string, error) {
	input = strings.TrimSpace(input)

	// Accept paths pasted with surrounding quotes.
	if len(input) >= 2 {
		if (input[0] == '"' && input[len(input)-1] == '"') ||
			(input[0] == '\'' && input[len(input)-1] == '\'') {
			input = input[1 : len(input)-1]
		}
	}

	// Empty input means the current terminal folder.
	if input == "" {
		input = "."
	}

	// Support the user's home folder on all platforms.
	if input == "~" ||
		strings.HasPrefix(input, "~/") ||
		strings.HasPrefix(input, "~\\") {

		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot find your home folder: %w", err)
		}

		if input == "~" {
			input = home
		} else {
			input = filepath.Join(home, input[2:])
		}
	}

	directory, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("invalid folder path: %w", err)
	}

	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return "", fmt.Errorf(
			"folder does not exist; create it first or choose another",
		)
	}
	if err != nil {
		return "", fmt.Errorf("cannot access this folder: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("choose a folder, not a file")
	}

	return directory, nil
}

func normalizeServiceName(input string) string {
	return strings.ToLower(strings.Join(strings.Fields(input), "-"))
}

func validateServiceName(input string) error {
	name := normalizeServiceName(input)

	// Leave room for the suffixes added to Kubernetes service names.
	if len(name) == 0 || len(name) > 46 {
		return fmt.Errorf("use a name between 1 and 46 characters")
	}

	if !serviceNamePattern.MatchString(name) {
		return fmt.Errorf(
			"start with a letter; use letters, numbers and hyphens; end with a letter or number",
		)
	}

	if reservedNamePattern.MatchString(name) {
		return fmt.Errorf("this name is reserved on Windows; choose another")
	}

	return nil
}
