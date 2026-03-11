package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/revotale/infrastructure/internal/routesgen"
)

func main() {
	var (
		check      = flag.Bool("check", false, "verify generated files are up to date")
		configPath = flag.String("config", "nginx/routes.yml", "path to routes manifest")
		readmePath = flag.String("readme", "nginx/README.md", "path to README")
	)

	flag.Parse()

	if err := run(*configPath, *readmePath, *check); err != nil {
		fmt.Fprintf(os.Stderr, "routesgen: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath string, readmePath string, check bool) error {
	cfg, err := routesgen.LoadConfig(configPath)
	if err != nil {
		return err
	}

	readmeBytes, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("read README: %w", err)
	}

	outputs, err := routesgen.GenerateOutputs(cfg, string(readmeBytes))
	if err != nil {
		return err
	}

	changes, err := syncOutputs(outputs, routesgen.GeneratedRouteFiles(cfg), check)
	if err != nil {
		return err
	}

	if check && len(changes) > 0 {
		sort.Strings(changes)
		return fmt.Errorf("generated files are out of date:\n%s", strings.Join(changes, "\n"))
	}

	return nil
}

func syncOutputs(outputs map[string]string, expectedRouteFiles []string, check bool) ([]string, error) {
	var changes []string

	for path, content := range outputs {
		current, err := os.ReadFile(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}

		if string(current) == content {
			continue
		}

		changes = append(changes, path)
		if check {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return nil, fmt.Errorf("write %s: %w", path, err)
		}
	}

	expected := make(map[string]struct{}, len(expectedRouteFiles))
	for _, path := range expectedRouteFiles {
		expected[path] = struct{}{}
	}

	actualRouteFiles, err := filepath.Glob("nginx/templates/routes/*.conf.template")
	if err != nil {
		return nil, fmt.Errorf("list generated route templates: %w", err)
	}

	for _, path := range actualRouteFiles {
		path = filepath.ToSlash(path)
		if _, ok := expected[path]; ok {
			continue
		}

		changes = append(changes, path)
		if check {
			continue
		}

		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("remove stale file %s: %w", path, err)
		}
	}

	return changes, nil
}
