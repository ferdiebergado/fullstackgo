package service

import "errors"

var ErrModelNotFound = errors.New("model not found")
var ErrModelExists = errors.New("model already exists")
