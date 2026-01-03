package user

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrPasswordIncorrect = errors.New("password incorrect")
