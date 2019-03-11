package file

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	contents       = "xxxx"
	contentsLonger = "yyyyyy"
)

func TestOverwrite(t *testing.T) {
	f, err := ioutil.TempFile("", "write_test")
	require.NoError(t, err)

	err = Overwrite(f.Name(), contents)
	require.NoError(t, err)

	b, err := ioutil.ReadFile(f.Name())
	require.Equal(t, contents, string(b))

	err = Overwrite(f.Name(), contentsLonger)
	require.NoError(t, err)

	b, err = ioutil.ReadFile(f.Name())
	require.Equal(t, contentsLonger, string(b))

	err = Overwrite(f.Name(), contents)
	require.NoError(t, err)

	b, err = ioutil.ReadFile(f.Name())
	require.Equal(t, contents, string(b))
}
