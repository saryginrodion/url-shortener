package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type Response[T any] struct {
	Ok   bool `json:"ok"`
	Data T    `json:"data"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Ok:      false,
		Message: message,
	}
}

func NewResponse[T any](val T) Response[T] {
	return Response[T]{
		Ok:   true,
		Data: val,
	}
}

func WriteJsonErrorResponse(w http.ResponseWriter, err error, status int) error {
	return WriteJsonResponse(
		w,
		NewErrorResponse(err.Error()),
		status,
	)
}

func WriteJsonResponse[T any](w http.ResponseWriter, val T, status int) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err = w.Write(bytes); err != nil {
		return err
	}

	return nil
}
