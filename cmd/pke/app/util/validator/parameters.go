package validator

import (
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/pkg/errors"
)

// NotEmpty gives error, if any of the given args is empty. Map key is returned in the error message.
func NotEmpty(args map[string]interface{}) error {
	for k, v := range args {
		switch arg := v.(type) {
		case string:
			if arg == "" {
				return errors.Wrapf(constants.ErrValidationFailed, "missing %s", k)
			}
		case int32:
			if arg <= 0 {
				return errors.Wrapf(constants.ErrValidationFailed, "missing %s", k)
			}
		}
	}
	return nil
}

// Empty gives error, if any of the given args is not empty. Map key is returned in the error message.
func Empty(args map[string]interface{}) error {
	for k, v := range args {
		switch arg := v.(type) {
		case string:
			if arg != "" {
				return errors.Wrapf(constants.ErrValidationFailed, "missing %s", k)
			}
		case int32:
			if arg > 0 {
				return errors.Wrapf(constants.ErrValidationFailed, "missing %s", k)
			}
		}
	}
	return nil
}
