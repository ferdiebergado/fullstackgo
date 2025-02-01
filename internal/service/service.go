package service

import "errors"

var ErrUserNotFound = errors.New("user does not exists")
var ErrEmailTaken = errors.New("email is already taken")
