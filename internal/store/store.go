package store

type Store interface {
	CreateShortURL(originalURL string) (string, error)
	ResolveShortURL(code string) (string, bool, error)
	Summary() (Summary, error)
	Top(limit int) ([]LinkInfo, error)
	Recent(limit int) ([]LinkInfo, error)
	Close() error
}
