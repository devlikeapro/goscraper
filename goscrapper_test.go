package goscraper

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestScrape function to test Scrape behavior.
func TestScrape(t *testing.T) {
	t.Run("valid URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<html><head><title>Test</title></head><body><h1>Test</h1></body></html>`))
		}))
		defer server.Close()

		_, err := Scrape(server.URL, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		_, err := Scrape("invalid-url", 5)
		if err == nil {
			t.Fatal("expected error for invalid URL, got nil")
		}
	})

	t.Run("max redirect handling", func(t *testing.T) {
		redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/", http.StatusFound)
		}))
		defer redirectServer.Close()

		_, err := Scrape(redirectServer.URL, 0) // No redirects allowed
		if err == nil {
			t.Fatal("expected error due to max redirect limit")
		}
	})
}

// TestScrapeIntegration function to test Scrape behavior.
// It's not unit test, it's integration test.
// The values might change over time.
func TestScrapeIntegration(t *testing.T) {
	t.Run("https://www.w3.org/", func(t *testing.T) {
		url := "https://www.w3.org/"
		expectedPreview := DocumentPreview{
			Icon:        "https://www.w3.org/assets/logos/w3c/favicon-180.png",
			Name:        "W3C",
			Title:       "W3C",
			Description: "The World Wide Web Consortium (W3C) develops standards and guidelines to help everyone build a web based on the principles of accessibility, internationalization, privacy and security.",
			Images:      []string{"https://www.w3.org/assets/website-2021/images/w3c-opengraph-image.png"},
			Link:        "https://www.w3.org/",
		}
		doc, err := Scrape(url, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc == nil {
			t.Fatalf("expected non-nil document")
		}
		assert.Equal(t, doc.Preview, expectedPreview)
	})

	t.Run("withContext - success", func(t *testing.T) {
		url := "https://www.w3.org/"
		ctx, _ := context.WithTimeout(t.Context(), 10*time.Second)
		_, err := ScrapeWithContext(ctx, url, 5)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("withContext - timeout", func(t *testing.T) {
		url := "https://www.w3.org/"
		ctx, _ := context.WithTimeout(t.Context(), time.Nanosecond)
		_, err := ScrapeWithContext(ctx, url, 5)
		if err == nil {
			t.Fatalf("expected error due to timeout")
		}
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})
}
