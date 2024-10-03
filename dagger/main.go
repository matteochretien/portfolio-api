// A generated module for Api functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"bufio"
	"context"
	"fmt"
	"strings"
)

type Api struct{}

func (a *Api) Migrate(ctx context.Context,
	env *File,
	dir *Directory,
	// +optional
	destination string,
) (string, error) {
	contents, err := env.Contents(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to read .env file: %w", err)
	}

	envs, err := ParseEnvFile(contents)
	if err != nil {
		return "", fmt.Errorf("failed to parse .env file: %w", err)
	}

	// Get the DSN from the .env file
	dsn, ok := envs["POSTGRES_DSN"]
	if !ok {
		return "", fmt.Errorf("POSTGRES_DSN not found in .env file")
	}

	cmd := []string{"tern", "migrate", "--migrations", "/migrations", "--conn-string", dsn}
	if destination != "" {
		cmd = append(cmd, "--destination", destination)
	}

	stdout, err := dag.Container().
		From("golang:1.22-alpine3.19").
		WithExec([]string{"go", "install", "github.com/jackc/tern/v2@latest"}).
		WithMountedDirectory("/migrations", dir).
		WithExec(cmd).
		Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run tern migrate: %w", err)
	}

	return stdout, nil
}

func (a *Api) New(ctx context.Context,
	dir *Directory,
	name string,
) (*File, error) {
	directory := dag.Container().
		From("golang:1.22-alpine3.19").
		WithExec([]string{"go", "install", "github.com/jackc/tern/v2@latest"}).
		WithMountedDirectory("/migrations", dir).
		WithExec([]string{"tern", "new", "--migrations", "/migrations", name}).
		Directory("/migrations")

	entries, err := directory.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory entries: %w", err)
	}

	for _, entry := range entries {
		if strings.Contains(entry, name) {
			return directory.File(entry), nil
		}
	}

	return nil, fmt.Errorf("migration file not found")
}

// ParseEnvFile parse the content of a .env file and returns a map of key-value pairs
func ParseEnvFile(content string) (map[string]string, error) {
	envMap := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()

		// Trim spaces and ignore comments and empty lines
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Ignore malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		if len(value) >= 2 {
			if (value[0] == '\'' && value[len(value)-1] == '\'') || (value[0] == '"' && value[len(value)-1] == '"') {
				value = value[1 : len(value)-1]
			}
		}

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}
