package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsolutise(t *testing.T) {
	testCases := []struct {
		base     string
		p        string
		expected string
	}{
		{"/", "/usr/local/bin/", "/usr/local/bin"},
		{"/", "./usr/local/bin/", "/usr/local/bin"},
		{"/tmp/", "/usr/local/bin/", "/usr/local/bin"},
		{"/tmp/", "./usr/local/bin/", "/tmp/usr/local/bin"},
	}

	for _, tc := range testCases {
		abs := Absolutise(tc.base, tc.p)
		assert.Equal(t, tc.expected, abs)
	}
}
