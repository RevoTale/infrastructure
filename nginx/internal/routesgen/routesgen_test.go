package routesgen

import (
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

func TestGenerateOutputsIncludesExpectedTemplates(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)

	outputs, err := GenerateOutputs(cfg)
	if err != nil {
		t.Fatalf("generate outputs: %v", err)
	}

	for _, name := range GeneratedRouteFilenames(cfg) {
		if _, ok := outputs[name]; !ok {
			t.Fatalf("missing generated output %s", name)
		}
	}

	if !strings.Contains(outputs["01-http-redirect.conf.template"], "return 301 https://$host$request_uri;") {
		t.Fatalf("expected HTTP redirect template to redirect to HTTPS")
	}

	if !strings.Contains(outputs["99-default.conf.template"], "return 404;") {
		t.Fatalf("expected default template to return fallback status")
	}
}

func TestGeneratedRouteFilenamesMatchExpectedOrder(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)

	got := GeneratedRouteFilenames(cfg)
	want := []string{
		"01-http-redirect.conf.template",
		"10-landing.conf.template",
		"11-app-redirect.conf.template",
		"12-console.conf.template",
		"13-status.conf.template",
		"99-default.conf.template",
	}

	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected generated filenames:\n%s", strings.Join(got, "\n"))
	}
}

func TestGeneratedTemplatesUseVariableProxyPass(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)

	outputs, err := GenerateOutputs(cfg)
	if err != nil {
		t.Fatalf("generate outputs: %v", err)
	}

	landing := outputs["10-landing.conf.template"]
	if !strings.Contains(landing, "set $upstream public-web:8080;") {
		t.Fatalf("expected landing template to set upstream via variable")
	}
	if !strings.Contains(landing, "proxy_pass http://$upstream;") {
		t.Fatalf("expected landing template to use variable proxy_pass")
	}
	if !strings.Contains(landing, "${PROJECT_DOMAIN}") {
		t.Fatalf("expected landing template to retain PROJECT_DOMAIN placeholder")
	}
}

func TestGeneratedTemplatesRenderRegexLocationsForNginx(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)

	outputs, err := GenerateOutputs(cfg)
	if err != nil {
		t.Fatalf("generate outputs: %v", err)
	}

	console := outputs["12-console.conf.template"]
	if !strings.Contains(console, "location ~ ^/downloads/(.*)\\.zip$ {") {
		t.Fatalf("expected regex route location to use nginx regex syntax")
	}
	if !strings.Contains(console, "rewrite ^/downloads/(.*)\\.zip$ /archives/$1.zip break;") {
		t.Fatalf("expected regex route location to include rewrite")
	}
	if !strings.Contains(console, "set $upstream archive-service:8082;") {
		t.Fatalf("expected regex route location to use masked upstream")
	}

	landing := outputs["10-landing.conf.template"]
	if !strings.Contains(landing, "location ~ ^/assets/cache/(thumb|hero)/(.*)$ {") {
		t.Fatalf("expected shared rewrite to render as nginx regex location")
	}
	if !strings.Contains(landing, "rewrite ^/assets/cache/(thumb|hero)/(.*)$ /unsafe/fit-in/$1/plain/http://asset-origin:8081/$2 break;") {
		t.Fatalf("expected shared rewrite to include masked rewrite target")
	}
}

func TestRenderedTemplatesWithProjectDomain(t *testing.T) {
	t.Parallel()

	cfg := mustLoadFixtureConfig(t)

	outputs, err := GenerateOutputs(cfg)
	if err != nil {
		t.Fatalf("generate outputs: %v", err)
	}

	var sawExampleDomain bool
	for _, name := range GeneratedRouteFilenames(cfg) {
		rendered := strings.ReplaceAll(outputs[name], projectDomainVar, "example.com")
		if strings.Contains(rendered, projectDomainVar) {
			t.Fatalf("rendered template still contains %s in %s", projectDomainVar, name)
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

	cfg, err := LoadConfig(fixturePath("testdata", "masked-routes.yml"))
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
