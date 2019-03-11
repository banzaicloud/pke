package runner

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	c := Cmd(ioutil.Discard, "echo", "ok")
	err := c.Run()
	require.NoError(t, err)
}

func TestOutput(t *testing.T) {
	c := Cmd(ioutil.Discard, "echo", "ok")
	out, err := c.Output()
	require.NoError(t, err)
	require.Equal(t, []byte("ok\n"), out)
}

func TestPipeOut(t *testing.T) {
	c := Cmd(ioutil.Discard, "echo", "ok")
	o, err := c.StdoutPipe()
	require.NoError(t, err)

	err = c.Start()
	require.NoError(t, err)

	out, err := ioutil.ReadAll(o)
	require.NoError(t, err)
	require.Equal(t, []byte("ok\n"), out)

	err = c.Wait()
	require.NoError(t, err)
}
