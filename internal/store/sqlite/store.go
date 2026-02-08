package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"

	"url-shortener/internal/store"
	"url-shortener/internal/util"
)

const (
	shortCodeLen = 6
	maxAttempts  = 8
)

type Store struct {
	db *sql.DB
}

var _ store.Store = (*Store)(nil)

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) CreateShortURL(originalURL string) (string, error) {
	stmt, err := s.db.Prepare(`INSERT INTO urls(code, url, created_at) VALUES(?, ?, ?)`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	for i := 0; i < maxAttempts; i++ {
		code, err := util.RandomCode(shortCodeLen)
		if err != nil {
			return "", err
		}

		_, err = stmt.Exec(code, originalURL, time.Now().Unix())
		if err == nil {
			return code, nil
		}

		if isConstraintError(err) {
			continue
		}

		return "", err
	}

	return "", fmt.Errorf("failed to generate unique short code after %d attempts", maxAttempts)
}

func (s *Store) ResolveShortURL(code string) (string, bool, error) {
	var url string
	row := s.db.QueryRow(`SELECT url FROM urls WHERE code = ?`, code)
	if err := row.Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}

	_, _ = s.db.Exec(`UPDATE urls SET clicks = clicks + 1 WHERE code = ?`, code)
	return url, true, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func isConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code == sqlite3.ErrConstraint
	}
	return false
}

func (s *Store) Summary() (store.Summary, error) {
	var summary store.Summary
	row := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(clicks), 0) FROM urls`)
	if err := row.Scan(&summary.TotalURLs, &summary.TotalClicks); err != nil {
		return store.Summary{}, err
	}
	return summary, nil
}

func (s *Store) Top(limit int) ([]store.LinkInfo, error) {
	if limit <= 0 {
		return []store.LinkInfo{}, nil
	}
	rows, err := s.db.Query(`SELECT code, url, clicks, created_at FROM urls ORDER BY clicks DESC, created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []store.LinkInfo
	for rows.Next() {
		var info store.LinkInfo
		if err := rows.Scan(&info.Code, &info.URL, &info.Clicks, &info.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, info)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Store) Recent(limit int) ([]store.LinkInfo, error) {
	if limit <= 0 {
		return []store.LinkInfo{}, nil
	}
	rows, err := s.db.Query(`SELECT code, url, clicks, created_at FROM urls ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []store.LinkInfo
	for rows.Next() {
		var info store.LinkInfo
		if err := rows.Scan(&info.Code, &info.URL, &info.Clicks, &info.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, info)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
