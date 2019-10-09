// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
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
