package analytics

import "time"

type ClicksByGeo struct {
	CountryCode string `json:"country_code"`
	Clicks      int    `json:"clicks"`
}

type ClicksByReferer struct {
	Referer string `json:"referer"`
	Clicks  int    `json:"clicks"`
}

type UrlStatistics struct {
	UrlID       string            `json:"url_id"`
	Date        time.Time         `json:"date"`
	TotalClicks int               `json:"total_clicks"`
	ByGeo       []ClicksByGeo     `json:"by_geo"`
	ByReferer   []ClicksByReferer `json:"by_referer"`
}
