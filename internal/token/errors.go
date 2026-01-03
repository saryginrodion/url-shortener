package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"

	"roadmap.restapi/internal/errormapper"
)

var (
	ErrTokenExpired             = errors.New("token expired")
	ErrTokenMalformed           = errors.New("token malformed")
	ErrTokenInvalid             = errors.New("token invalid")
	ErrTokenSignatureInvalid    = errors.New("token signature is invalid")
	ErrRefreshNotWhitelisted    = errors.New("refresh token is not in whitelist")
)

var jwtErrMapper = errormapper.NewErrorMapper(
	errormapper.NewMapping(jwt.ErrTokenMalformed, ErrTokenMalformed),
	errormapper.NewMapping(jwt.ErrTokenExpired, ErrTokenExpired),
	errormapper.NewMapping(jwt.ErrTokenMalformed, ErrTokenInvalid),
	errormapper.NewMapping(jwt.ErrTokenSignatureInvalid, ErrTokenSignatureInvalid),

	errormapper.NewMapping(jwt.ErrTokenUnverifiable, ErrTokenInvalid),
)

func TranslateJWTError(err error) error {
	if err == nil {
		return nil
	}

	mappedErr, matched := jwtErrMapper.MapAndCheck(err)
	if matched {
		return mappedErr
	}

	return ErrTokenInvalid
}
