package store

type LinkInfo struct {
	Code      string `json:"code"`
	URL       string `json:"url"`
	Clicks    int64  `json:"clicks"`
	CreatedAt int64  `json:"created_at"`
}

type Summary struct {
	TotalURLs   int64 `json:"total_urls"`
	TotalClicks int64 `json:"total_clicks"`
}
