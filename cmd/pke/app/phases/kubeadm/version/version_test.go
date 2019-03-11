package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidVersion(t *testing.T) {
	testCases := []struct {
		version string
		valid   bool
	}{
		{"0.0.1", false},
		{"1.12.0", true},
		{"1.12.1", true},
		{"1.12.2", true},
		{"1.12.3", true},
		{"1.12.4", true},
		{"1.12.5", true},
		{"1.12.6", false},
		{"1.13.0", true},
		{"1.13.1", true},
		{"1.13.2", true},
		{"1.13.3", true},
		{"1.13.4", true},
		{"1.13.5", false},
	}

	for _, tc := range testCases {
		err := validVersion(tc.version, constraint)
		if !tc.valid {
			assert.Error(t, err, tc.version)
		} else {
			assert.NoError(t, err, tc.version)
		}
	}
}
