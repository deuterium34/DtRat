package errs

import "errors"

var (
	ErrUnsupportedOs = errors.New("Эта ОС не поддерживается")
)
