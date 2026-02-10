package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/StealthBadger747/ShortSlug/internal/bot"
	"github.com/StealthBadger747/ShortSlug/internal/store"
)

type Server struct {
	frontendDir string
	store       store.Store
	capVerifier *bot.CapVerifier
	capEndpoint string
	password    string
}

func New(frontendDir string, store store.Store, capVerifier *bot.CapVerifier, capEndpoint string, password string) *Server {
	return &Server{
		frontendDir: frontendDir,
		store:       store,
		capVerifier: capVerifier,
		capEndpoint: capEndpoint,
		password:    password,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/api/shorten_url" {
		s.handleShorten(w, r)
		return
	}

	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/analytics/") {
		s.handleAnalytics(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		s.serveIndex(w, r)
		return
	}

	if s.tryServeStatic(w, r) {
		return
	}

	s.handleRedirect(w, r)
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, r, http.StatusBadRequest, "Invalid form data.")
		return
	}

	if s.password != "" {
		if r.FormValue("password") != s.password {
			writeError(w, r, http.StatusUnauthorized, "Invalid password.")
			return
		}
	}

	if s.capVerifier != nil && s.capVerifier.Enabled() {
		token := r.FormValue("cap-token")
		if err := s.capVerifier.Verify(r.Context(), token); err != nil {
			writeError(w, r, http.StatusBadRequest, "Bot verification failed.")
			return
		}
	}

	originalURL := strings.TrimSpace(r.FormValue("url"))
	if originalURL == "" {
		writeError(w, r, http.StatusBadRequest, "Please enter a URL before shortening.")
		return
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	parsed, err := url.ParseRequestURI(originalURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		writeError(w, r, http.StatusBadRequest, "That URL doesn't look valid. Check the format and try again.")
		return
	}

	code, err := s.store.CreateShortURL(originalURL)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "Failed to create short URL.")
		return
	}

	baseURL := fmt.Sprintf("%s://%s", schemeForRequest(r), r.Host)
	shortURL := baseURL + "/" + code

	if isHtmxRequest(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, renderHtmxResult(shortURL))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":         "200",
		"status_message": "OK",
		"short_url":      shortURL,
	})
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, ok, err := s.store.ResolveShortURL(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, "404 NOT FOUND!")
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (s *Server) handleAnalytics(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/analytics/summary":
		summary, err := s.store.Summary()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, summary)
	case "/api/analytics/top":
		limit := parseLimit(r.URL.Query().Get("limit"))
		links, err := s.store.Top(limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, links)
	case "/api/analytics/recent":
		limit := parseLimit(r.URL.Query().Get("limit"))
		links, err := s.store.Recent(limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, links)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) tryServeStatic(w http.ResponseWriter, r *http.Request) bool {
	cleanPath := path.Clean(r.URL.Path)
	if strings.Contains(cleanPath, "..") {
		w.WriteHeader(http.StatusBadRequest)
		return true
	}

	fsPath := filepath.Join(s.frontendDir, filepath.FromSlash(strings.TrimPrefix(cleanPath, "/")))
	absPath, err := filepath.Abs(fsPath)
	if err != nil {
		return false
	}

	if !strings.HasPrefix(absPath, s.frontendDir) {
		w.WriteHeader(http.StatusBadRequest)
		return true
	}

	info, err := os.Stat(absPath)
	if err != nil || info.IsDir() {
		return false
	}

	http.ServeFile(w, r, absPath)
	return true
}

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request, file string) {
	filePath := filepath.Join(s.frontendDir, file)
	http.ServeFile(w, r, filePath)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(s.frontendDir, "index.html")
	tmpl, err := template.ParseFiles(indexPath)
	if err != nil {
		http.ServeFile(w, r, indexPath)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, struct {
		CapAPIEndpoint  string
		PasswordEnabled bool
	}{
		CapAPIEndpoint:  s.capEndpoint,
		PasswordEnabled: s.password != "",
	})
}

func isHtmxRequest(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("HX-Request"), "true")
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	if isHtmxRequest(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		escaped := template.HTMLEscapeString(message)
		_, _ = io.WriteString(w, "<div class=\"alert error\">"+escaped+"</div>")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":         fmt.Sprintf("%d", status),
		"status_message": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func renderHtmxResult(shortURL string) string {
	escaped := template.HTMLEscapeString(shortURL)
	return "<div class=\"result\">" +
		"<p class=\"result-label\">Short URL</p>" +
		"<a class=\"result-link\" href=\"" + escaped + "\" target=\"_blank\" rel=\"noopener noreferrer\">" +
		escaped +
		"</a>" +
		"</div>"
}

func schemeForRequest(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	return "http"
}

func parseLimit(raw string) int {
	if raw == "" {
		return 10
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return 10
	}
	if val > 100 {
		return 100
	}
	return val
}
