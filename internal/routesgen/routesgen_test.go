package routesgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseConfigValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name: "duplicate ids",
			input: `
version: 1
http:
  redirect_to_https: true
fallback:
  status: 404
routes:
  - id: dup
    host: apex
    locations:
      - match_type: prefix
        match: /
        upstream: cms:3000
  - id: dup
    host: blog
    locations:
      - match_type: prefix
        match: /
        upstream: blog:8080
`,
			wantErr: `id "dup" is defined more than once`,
		},
		{
			name: "invalid host token",
			input: `
version: 1
http:
  redirect_to_https: true
fallback:
  status: 404
routes:
  - id: bad-host
    host: bad_host
    locations:
      - match_type: prefix
        match: /
        upstream: cms:3000
`,
			wantErr: `invalid host "bad_host"`,
		},
		{
			name: "invalid redirect shape",
			input: `
version: 1
http:
  redirect_to_https: true
fallback:
  status: 404
routes:
  - id: redirect
    host: www
    redirect:
      to_host: apex
      status: 301
    locations:
      - match_type: prefix
        match: /
        upstream: cms:3000
`,
			wantErr: `must define exactly one of redirect or locations`,
		},
		{
			name: "invalid location type",
			input: `
version: 1
http:
  redirect_to_https: true
fallback:
  status: 404
routes:
  - id: root
    host: apex
    locations:
      - match_type: contains
        match: /
        upstream: cms:3000
`,
			wantErr: `unsupported match_type "contains"`,
		},
		{
			name: "invalid regex rewrite pair",
			input: `
version: 1
http:
  redirect_to_https: true
fallback:
  status: 404
routes:
  - id: root
    host: apex
    locations:
      - match_type: prefix
        match: /
        rewrite: /elsewhere
        upstream: cms:3000
`,
			wantErr: `rewrite is only supported for regex matches`,
		},
		{
			name: "conflicting host match precedence",
			input: `
version: 1
http:
  redirect_to_https: true
fallback:
  status: 404
routes:
  - id: root
    host: apex
    locations:
      - match_type: prefix
        match: /
        upstream: cms:3000
      - match_type: prefix
        match: /
        upstream: blog:8080
`,
			wantErr: `repeats prefix match "/"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ParseConfig([]byte(tt.input))
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGeneratedRouteTemplatesMatchRepo(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)
	readmePath := fixturePath("nginx", "README.md")
	readmeBytes, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	outputs, err := GenerateOutputs(cfg, string(readmeBytes))
	if err != nil {
		t.Fatalf("generate outputs: %v", err)
	}

	for _, path := range GeneratedRouteFiles(cfg) {
		expectedBytes, err := os.ReadFile(filepath.Clean(fixturePath(path)))
		if err != nil {
			t.Fatalf("read generated template %s: %v", path, err)
		}

		if got := outputs[path]; got != string(expectedBytes) {
			t.Fatalf("generated template mismatch for %s", path)
		}
	}
}

func TestGeneratedReadmeMatchesRepo(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)
	readmePath := fixturePath("nginx", "README.md")
	readmeBytes, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	outputs, err := GenerateOutputs(cfg, string(readmeBytes))
	if err != nil {
		t.Fatalf("generate outputs: %v", err)
	}

	if got := outputs[filepath.ToSlash(filepath.Join("nginx", "README.md"))]; got != string(readmeBytes) {
		t.Fatalf("generated README does not match repository copy")
	}
}

func TestRenderedTemplatesWithProjectDomain(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)
	readmeBytes, err := os.ReadFile(fixturePath("nginx", "README.md"))
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	outputs, err := GenerateOutputs(cfg, string(readmeBytes))
	if err != nil {
		t.Fatalf("generate outputs: %v", err)
	}

	var sawExampleDomain bool
	for _, path := range GeneratedRouteFiles(cfg) {
		rendered := strings.ReplaceAll(outputs[path], projectDomainVar, "example.com")
		if strings.Contains(rendered, projectDomainVar) {
			t.Fatalf("rendered template still contains %s in %s", projectDomainVar, path)
		}
		if strings.Contains(rendered, "example.com") {
			sawExampleDomain = true
		}
	}

	if !sawExampleDomain {
		t.Fatalf("expected rendered templates to include example.com")
	}
}

func mustLoadFixtureConfig(t *testing.T) Config {
	t.Helper()

	cfg, err := LoadConfig(fixturePath("nginx", "routes.yml"))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	return cfg
}

func fixturePath(parts ...string) string {
	base := filepath.Join("..", "..")
	all := append([]string{base}, parts...)
	return filepath.Join(all...)
}
