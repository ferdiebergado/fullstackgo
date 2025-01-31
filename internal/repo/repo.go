package repo

import (
	"errors"
)

var ErrNullValue = errors.New("not null constraint violation")
