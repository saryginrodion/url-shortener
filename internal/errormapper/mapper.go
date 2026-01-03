package errormapper

import (
	"errors"
	"log/slog"
)

type mapping struct {
	predicate func(err error) bool
	to        error
}

type ErrorMapper struct {
	mappings     []mapping
}

func NewPredicateMapping(predicate func(err error) bool, to error) mapping {
	return mapping{
		predicate: predicate,
		to:        to,
	}
}

func NewMapping(from error, to error) mapping {
	return NewPredicateMapping(
		func(err error) bool { return errors.Is(err, from) },
		to,
	)
}

func NewErrorMapper(mappings ...mapping) *ErrorMapper {
	return &ErrorMapper{
		mappings: mappings,
	}
}

func (em *ErrorMapper) MapAndLogUnmatched(err error, log *slog.Logger) error {
	mapped, matched := em.MapAndCheck(err)
	if !matched && log != nil {
		log.Error("unmatched err", "err", err)
	}

	return mapped
}

func (em *ErrorMapper) Map(err error) error {
	if err == nil {
		return nil
	}

	mapped, _ := em.MapAndCheck(err)
	return mapped
}

func (em *ErrorMapper) MapAndCheck(err error) (error, bool) {
	for _, mapping := range em.mappings {
		if mapping.predicate(err) {
			return mapping.to, true
		}
	}

	return err, false
}
