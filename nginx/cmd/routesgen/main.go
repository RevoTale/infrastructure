package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/revotale/infrastructure/nginx/internal/routesgen"
)

func main() {
	var (
		check     = flag.Bool("check", false, "verify generated files are up to date")
		clean     = flag.Bool("clean", false, "remove the output directory before writing generated files")
		config    = flag.String("config", "routes.yml", "path to routes manifest")
		outputDir = flag.String("output-dir", "templates/routes", "directory for generated route templates")
	)

	flag.Parse()

	if err := run(*config, *outputDir, *check, *clean); err != nil {
		fmt.Fprintf(os.Stderr, "routesgen: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath string, outputDir string, check bool, clean bool) error {
	cfg, err := routesgen.LoadConfig(configPath)
	if err != nil {
		return err
	}

	outputs, err := routesgen.GenerateOutputs(cfg)
	if err != nil {
		return err
	}

	changes, err := syncOutputs(outputDir, outputs, check, clean)
	if err != nil {
		return err
	}

	if check && len(changes) > 0 {
		sort.Strings(changes)
		return fmt.Errorf("generated files are out of date:\n%s", strings.Join(changes, "\n"))
	}

	return nil
}

func syncOutputs(outputDir string, outputs map[string]string, check bool, clean bool) ([]string, error) {
	if clean && !check {
		if err := os.RemoveAll(outputDir); err != nil {
			return nil, fmt.Errorf("remove output dir %s: %w", outputDir, err)
		}
	}

	if !check {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", outputDir, err)
		}
	}

	var changes []string
	for name, content := range outputs {
		path := filepath.Join(outputDir, name)

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

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return nil, fmt.Errorf("write %s: %w", path, err)
		}
	}

	actualFiles, err := filepath.Glob(filepath.Join(outputDir, "*.conf.template"))
	if err != nil {
		return nil, fmt.Errorf("list generated route templates: %w", err)
	}

	expectedFiles := make(map[string]struct{}, len(outputs))
	for name := range outputs {
		expectedFiles[filepath.Join(outputDir, name)] = struct{}{}
	}

	for _, path := range actualFiles {
		if _, ok := expectedFiles[path]; ok {
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
