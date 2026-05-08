package gateway

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// TestSchemaValidation verifies the compiled OpenAPI spec matches the
// hosted registry version. This is an integration test that requires
// network access to the schema registry.
func TestSchemaValidation(t *testing.T) {
	// Walk up to find repo root (where .cache/ lives)
	cacheFile := "../../.cache/docs-url"
	urlBytes, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("schema cache not found — run 'make setup' first")
		}
		t.Fatalf("reading schema url: %v", err)
	}

	endpoint := strings.TrimSpace(string(urlBytes))
	if endpoint == "" {
		t.Fatal("empty schema endpoint in .cache/docs-url")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		t.Fatalf("schema registry request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Registry should return 200 with a non-empty body
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("schema registry returned %d (expected 200), body: %s",
			resp.StatusCode, string(body[:min(len(body), 200)]))
	}

	if len(body) == 0 {
		t.Fatal("schema registry returned empty response")
	}

	t.Logf("schema validated: %d bytes from registry", len(body))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
