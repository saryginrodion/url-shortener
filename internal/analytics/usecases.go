package analytics

import (
	"context"
)

type UseCases struct {
	repo URLStatisticsRepository
}

func NewUseCases(repo URLStatisticsRepository) *UseCases {
	return &UseCases{
		repo: repo,
	}
}

func (s *UseCases) AddClick(
	ctx context.Context,
	urlID string,
	countryCode *string,
	referer *string,
) error {
	if err := s.repo.AddClick(ctx, urlID); err != nil {
		return err
	}

	if countryCode != nil {
		if err := s.repo.AddGeoClick(ctx, urlID, *countryCode); err != nil {
			return err
		}
	}

	if referer != nil {
		if err := s.repo.AddRefererClick(ctx, urlID, *referer); err != nil {
			return err
		}
	}
	return nil
}
