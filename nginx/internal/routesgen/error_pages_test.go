package routesgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStaticErrorPagesExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		contains []string
	}{
		{
			name: "404 page",
			path: fixturePath("static", "__error_pages", "404.html"),
			contains: []string{
				"RevoTale Edge",
				"404",
				"Return Home",
			},
		},
		{
			name: "429 page",
			path: fixturePath("static", "__error_pages", "429.html"),
			contains: []string{
				"RevoTale Edge",
				"429",
				"Too Many Requests",
			},
		},
		{
			name: "502 page",
			path: fixturePath("static", "__error_pages", "502.html"),
			contains: []string{
				"RevoTale Edge",
				"502",
				"Bad Gateway",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(filepath.Clean(tt.path))
			if err != nil {
				t.Fatalf("read static error page: %v", err)
			}

			content := string(data)
			for _, marker := range tt.contains {
				if !strings.Contains(content, marker) {
					t.Fatalf("expected %s to contain %q", tt.path, marker)
				}
			}
		})
	}
}

func TestServerCommonReferencesStaticErrorPages(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(filepath.Clean(fixturePath("includes", "server-common.conf")))
	if err != nil {
		t.Fatalf("read server-common.conf: %v", err)
	}

	content := string(data)
	for _, marker := range []string{
		"error_page 404 /__error_pages/404.html;",
		"error_page 429 /__error_pages/429.html;",
		"error_page 500 /__error_pages/500.html;",
		"error_page 502 /__error_pages/502.html;",
		"error_page 503 /__error_pages/503.html;",
		"error_page 504 /__error_pages/504.html;",
		"error_page 501 505 507 /__error_pages/5xx.html;",
		"location ^~ /__error_pages/ {",
		"root /usr/share/nginx/html;",
	} {
		if !strings.Contains(content, marker) {
			t.Fatalf("expected server-common.conf to contain %q", marker)
		}
	}
}
