package url

import "errors"

var (
	ErrURLNotFound      = errors.New("url not found")
	ErrURLAlreadyExists = errors.New("url already exists")
	ErrUserIsNotAuthor  = errors.New("this user is not author of url")
)
