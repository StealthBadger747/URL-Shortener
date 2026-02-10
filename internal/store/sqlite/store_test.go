package sqlite

import (
	"path/filepath"
	"testing"
)

func TestStoreCreateResolveAndAnalytics(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	code, err := store.CreateShortURL("http://example.com")
	if err != nil {
		t.Fatalf("create short url: %v", err)
	}
	code2, err := store.CreateShortURL("http://example.com")
	if err != nil {
		t.Fatalf("create short url again: %v", err)
	}
	if code2 != code {
		t.Fatalf("expected same code for same url, got %s vs %s", code, code2)
	}

	url, ok, err := store.ResolveShortURL(code)
	if err != nil {
		t.Fatalf("resolve short url: %v", err)
	}
	if !ok {
		t.Fatalf("expected url to resolve")
	}
	if url != "http://example.com" {
		t.Fatalf("unexpected url: %s", url)
	}

	summary, err := store.Summary()
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.TotalURLs != 1 {
		t.Fatalf("expected 1 url, got %d", summary.TotalURLs)
	}
	if summary.TotalClicks != 1 {
		t.Fatalf("expected 1 click, got %d", summary.TotalClicks)
	}

	top, err := store.Top(5)
	if err != nil {
		t.Fatalf("top: %v", err)
	}
	if len(top) != 1 || top[0].Code != code {
		t.Fatalf("unexpected top results")
	}

	recent, err := store.Recent(5)
	if err != nil {
		t.Fatalf("recent: %v", err)
	}
	if len(recent) != 1 || recent[0].Code != code {
		t.Fatalf("unexpected recent results")
	}
}
