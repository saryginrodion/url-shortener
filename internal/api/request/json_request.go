package request

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

func ParseAndValidateJson[T any](v *validator.Validate, b io.ReadCloser, res T) (*T, error) {
	val, err := io.ReadAll(b)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &res)

	if err != nil {
		return nil, err
	}

	err = v.Struct(res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

