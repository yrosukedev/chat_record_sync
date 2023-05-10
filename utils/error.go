package utils

import (
	"errors"
	"strings"
)

func CombineErrors[E error](errs []E) error {
	var msgs []string
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return errors.New(strings.Join(msgs, "\n"))
}
