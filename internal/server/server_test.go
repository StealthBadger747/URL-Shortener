package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"url-shortener/internal/store/sqlite"
)

func TestShortenRedirectAndAnalytics(t *testing.T) {
	frontendDir := t.TempDir()
	if err := writeIndex(frontendDir); err != nil {
		t.Fatalf("write index: %v", err)
	}

	store, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	srv := httptest.NewServer(New(frontendDir, store, nil, "", ""))
	defer srv.Close()

	form := url.Values{}
	form.Set("url", "example.com")
	resp, err := http.PostForm(srv.URL+"/api/shorten_url", form)
	if err != nil {
		t.Fatalf("post form: %v", err)
	}
	defer resp.Body.Close()

	var body struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if body.ShortURL == "" {
		t.Fatalf("expected short url")
	}

	resp2, err := http.PostForm(srv.URL+"/api/shorten_url", form)
	if err != nil {
		t.Fatalf("post form again: %v", err)
	}
	defer resp2.Body.Close()

	var body2 struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&body2); err != nil {
		t.Fatalf("decode json 2: %v", err)
	}
	if body2.ShortURL != body.ShortURL {
		t.Fatalf("expected same short url for same input")
	}

	parts := strings.Split(strings.TrimSuffix(body.ShortURL, "/"), "/")
	code := parts[len(parts)-1]

	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	redirResp, err := client.Get(srv.URL + "/" + code)
	if err != nil {
		t.Fatalf("get redirect: %v", err)
	}
	defer redirResp.Body.Close()

	if redirResp.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", redirResp.StatusCode)
	}

	summaryResp, err := http.Get(srv.URL + "/api/analytics/summary")
	if err != nil {
		t.Fatalf("get summary: %v", err)
	}
	defer summaryResp.Body.Close()

	var summary struct {
		TotalURLs   int64 `json:"total_urls"`
		TotalClicks int64 `json:"total_clicks"`
	}
	if err := json.NewDecoder(summaryResp.Body).Decode(&summary); err != nil {
		t.Fatalf("decode summary: %v", err)
	}
	if summary.TotalURLs != 1 || summary.TotalClicks != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	topResp, err := http.Get(srv.URL + "/api/analytics/top?limit=5")
	if err != nil {
		t.Fatalf("get top: %v", err)
	}
	defer topResp.Body.Close()

	var top []map[string]any
	if err := json.NewDecoder(topResp.Body).Decode(&top); err != nil {
		t.Fatalf("decode top: %v", err)
	}
	if len(top) != 1 {
		t.Fatalf("expected 1 top entry, got %d", len(top))
	}
}

func writeIndex(dir string) error {
	return os.WriteFile(filepath.Join(dir, "index.html"), []byte("ok"), 0644)
}
