# ShortSlug
## Architecture
- HTMX frontend served as static HTML/CSS.
- Go stdlib HTTP server for the API and static assets.
- SQLite for persistence.
  - Short codes are random, but the Go version reuses the same code for identical long URLs (store-and-reuse).

## Note
This project is also hosted on my server in my apartment.
I went a bit overboard with this implementation than was probably expected, but I have been wanting to create ShortSlug for myself for a while now and this was just a good opportunity. My server environment is exclusively in Docker so it didn't take too long to set up.

## How to run/build
The Go implementation uses SQLite for storage and serves the static HTMX frontend.

Run:
```bash
go run ./cmd/shortslug
```

Docker:
```bash
./run_docker.sh
```

Environment variables:
 - `SERVER_PORT` (default `8080`)
 - `FRONTEND_DIR` (default `static` if present)
 - `DATABASE_PATH` (default `shortslug.db`)
 - `CAP_SITEVERIFY_URL` (Cap siteverify endpoint; enables bot filtering)
 - `CAP_SECRET` (Cap secret key)
 - `CAP_API_ENDPOINT` (Cap widget API endpoint, used in `static/index.html`)
 - `SHORTEN_PASSWORD` (optional; if set, requires matching password to shorten)

Analytics endpoints (JSON):
 - `GET /api/analytics/summary`
 - `GET /api/analytics/top?limit=10`
 - `GET /api/analytics/recent?limit=10`

Bot filtering (Cap):
 - Include `CAP_SITEVERIFY_URL`, `CAP_SECRET`, and `CAP_API_ENDPOINT` to enable.

Database migrations:
 - Managed by `goose` and embedded in the binary.
 - Migrations live in `internal/store/sqlite/migrations`.

## Sources/Third Party Libraries:
- htmx for the frontend
- For the chain favicon: https://www.favicon-generator.org/search/---/Chain
