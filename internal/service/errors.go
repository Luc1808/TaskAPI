package service

import (
	"errors"
	"fmt"
)

var (
	ErrValidation   = errors.New("validation error")
	ErrDataNotFound = errors.New("not found")
)

func WrapValidation(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrValidation, err.Error())
}
