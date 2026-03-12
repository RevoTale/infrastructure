package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWithCleanRemovesStaleFiles(t *testing.T) {
	t.Parallel()

	outputDir := t.TempDir()
	stalePath := filepath.Join(outputDir, "stale.conf.template")
	if err := os.WriteFile(stalePath, []byte("stale"), 0o644); err != nil {
		t.Fatalf("seed stale file: %v", err)
	}

	if err := run(filepath.Join("..", "..", "testdata", "masked-routes.yml"), outputDir, false, true); err != nil {
		t.Fatalf("run routesgen: %v", err)
	}

	if _, err := os.Stat(stalePath); !os.IsNotExist(err) {
		t.Fatalf("expected stale file to be removed, got %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "10-landing.conf.template")); err != nil {
		t.Fatalf("expected generated route template: %v", err)
	}
}

func TestRunCheckDetectsDrift(t *testing.T) {
	t.Parallel()

	outputDir := t.TempDir()

	err := run(filepath.Join("..", "..", "testdata", "masked-routes.yml"), outputDir, true, false)
	if err == nil {
		t.Fatalf("expected check mode to detect missing generated files")
	}
	if !strings.Contains(err.Error(), "generated files are out of date") {
		t.Fatalf("unexpected error: %v", err)
	}
}
