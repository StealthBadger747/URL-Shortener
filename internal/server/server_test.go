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

	"github.com/StealthBadger747/ShortSlug/internal/store/sqlite"
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

	srv := httptest.NewServer(New(frontendDir, store, nil, "", "", "", "ShortSlug", "secret"))
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

	summaryReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/analytics/summary", nil)
	summaryReq.Header.Set("X-Analytics-Password", "secret")
	summaryResp, err := http.DefaultClient.Do(summaryReq)
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

	topReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/analytics/top?limit=5", nil)
	topReq.Header.Set("X-Analytics-Password", "secret")
	topResp, err := http.DefaultClient.Do(topReq)
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

func TestShortenUsesConfiguredPublicBaseURL(t *testing.T) {
	frontendDir := t.TempDir()
	if err := writeIndex(frontendDir); err != nil {
		t.Fatalf("write index: %v", err)
	}

	store, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	srv := httptest.NewServer(New(frontendDir, store, nil, "", "https://sho.rt", "", "ShortSlug", ""))
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

	if !strings.HasPrefix(body.ShortURL, "https://sho.rt/") {
		t.Fatalf("expected short URL to use configured base URL, got %q", body.ShortURL)
	}
}

func TestShortenUsesForwardedHeadersForBaseURL(t *testing.T) {
	frontendDir := t.TempDir()
	if err := writeIndex(frontendDir); err != nil {
		t.Fatalf("write index: %v", err)
	}

	store, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	h := New(frontendDir, store, nil, "", "", "", "ShortSlug", "")
	form := strings.NewReader("url=example.com")
	req := httptest.NewRequest(http.MethodPost, "/api/shorten_url", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "slug.example.com")
	req.Host = "internal:8080"
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	var body struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if !strings.HasPrefix(body.ShortURL, "https://slug.example.com/") {
		t.Fatalf("expected forwarded host/proto short URL, got %q", body.ShortURL)
	}
}

func TestShortenUsesRFCForwardedHeader(t *testing.T) {
	frontendDir := t.TempDir()
	if err := writeIndex(frontendDir); err != nil {
		t.Fatalf("write index: %v", err)
	}

	store, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	h := New(frontendDir, store, nil, "", "", "", "ShortSlug", "")
	form := strings.NewReader("url=example.com")
	req := httptest.NewRequest(http.MethodPost, "/api/shorten_url", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Forwarded", "for=1.2.3.4;proto=https;host=go.example.net")
	req.Host = "internal:8080"
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	var body struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if !strings.HasPrefix(body.ShortURL, "https://go.example.net/") {
		t.Fatalf("expected forwarded header short URL, got %q", body.ShortURL)
	}
}

func TestSecurityHeadersPresent(t *testing.T) {
	frontendDir := t.TempDir()
	if err := writeIndex(frontendDir); err != nil {
		t.Fatalf("write index: %v", err)
	}

	store, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	h := New(frontendDir, store, nil, "", "", "", "ShortSlug", "")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing nosniff header")
	}
	if rr.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("missing frame options header")
	}
	if rr.Header().Get("Referrer-Policy") != "no-referrer" {
		t.Fatalf("missing referrer policy header")
	}
	if rr.Header().Get("Content-Security-Policy") == "" {
		t.Fatalf("missing csp header")
	}
}

func writeIndex(dir string) error {
	return os.WriteFile(filepath.Join(dir, "index.html"), []byte("ok"), 0644)
}
