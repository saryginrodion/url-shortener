package analytics

import "context"

type URLStatisticsRepository interface {
	AddClick(ctx context.Context, urlID string) error
	AddRefererClick(ctx context.Context, urlID string, referer string) error
	AddGeoClick(ctx context.Context, urlID string, countryCode string) error
	Stats(ctx context.Context, urlID string) ([]UrlStatistics, error)
}
