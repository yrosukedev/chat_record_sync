package http_controller

import (
	"errors"
	"strings"
)

func combineErrors[E error](errs []E) error {
	var msgs []string
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return errors.New(strings.Join(msgs, "\n"))
}
