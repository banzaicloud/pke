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

package runner

import (
	"path/filepath"

	"emperror.dev/errors"

	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	name         string
	arg          []string
	w            io.Writer
	ts           time.Time
	errorMatcher func(string) bool
	*exec.Cmd
}

func trivialErrorMatcher(text string) bool {
	return strings.Contains(strings.ToLower(text), "error")
}

func Cmd(w io.Writer, name string, arg ...string) *Command {
	return &Command{
		name:         name,
		arg:          arg,
		w:            w,
		errorMatcher: trivialErrorMatcher,
		Cmd:          exec.Command(name, arg...),
	}
}

func (c *Command) ErrorMatcher(e func(string) bool) {
	c.errorMatcher = e
}

func (c *Command) CombinedOutput() ([]byte, error) {
	c.ts = time.Now()
	out, err := c.Cmd.CombinedOutput()
	_, _ = fmt.Fprintf(c.w, "%s %s err: %v %s\n", c.name, c.arg, err, time.Now().Sub(c.ts))
	if len(out) > 0 {
		_, _ = fmt.Fprintln(c.w, string(out))
	}
	return out, err
}

func (c *Command) CombinedOutputAsync() (string, error) {
	lastLine := ""
	firstError := ""

	c.ts = time.Now()

	stdout, err := c.Cmd.StdoutPipe()
	if err != nil {
		return lastLine, err
	}

	stdOutChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdOutChan <- scanner.Text()
		}
		close(stdOutChan)
	}()

	stderr, err := c.Cmd.StderrPipe()
	if err != nil {
		return lastLine, err
	}

	stdErrChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stdErrChan <- scanner.Text()
		}
		close(stdErrChan)
	}()

	err = c.Start()
	if err != nil {
		return lastLine, err
	}

	for stdErrChan != nil || stdOutChan != nil {
		var text string
		var more bool
		select {
		case text, more = <-stdOutChan:
			_, _ = fmt.Fprintln(c.w, "out>", text)
			if !more {
				stdOutChan = nil
			}
		case text, more = <-stdErrChan:
			_, _ = fmt.Fprintln(c.w, "err>", text)
			if !more {
				stdErrChan = nil
			}
		}

		if firstError == "" && c.errorMatcher != nil && c.errorMatcher(text) {
			firstError = text
		}
		lastLine = text
	}

	err = c.Wait()

	var target error = &exec.ExitError{}
	if errors.As(err, &target) {
		err = errors.WrapIff(err, "%s failed [%s]", filepath.Base(c.Args[0]), firstError)
	}

	if firstError == "" {
		firstError = lastLine
	}
	return firstError, err
}

func (c *Command) Output() ([]byte, error) {
	c.ts = time.Now()

	// Capture error output
	var stderr bytes.Buffer
	c.Cmd.Stderr = &stderr

	_, _ = fmt.Fprintf(c.w, "%s %s\n", c.name, c.arg)
	out, err := c.Cmd.Output()
	if len(out) > 0 {
		_, _ = fmt.Fprintf(c.w, "  out> %s\n", strings.ReplaceAll(string(out), "\n", "\n  out> "))
	}
	if stderr.Len() > 0 {
		_, _ = fmt.Fprintf(c.w, "  err> %s\n", strings.ReplaceAll(stderr.String(), "\n", "\n  err> "))
	}
	_, _ = fmt.Fprintf(c.w, "%s %s err: %v %s\n", c.name, c.arg, err, time.Now().Sub(c.ts))

	return out, err
}

func (c *Command) Run() error {
	c.ts = time.Now()
	err := c.Cmd.Run()
	_, _ = fmt.Fprintf(c.w, "%s %s err: %v %s\n", c.name, c.arg, err, time.Now().Sub(c.ts))
	return err
}

func (c *Command) Start() error {
	c.ts = time.Now()
	_, _ = fmt.Fprintf(c.w, "%s %s\n", c.name, c.arg)
	return c.Cmd.Start()
}

func (c *Command) Wait() error {
	err := c.Cmd.Wait()
	_, _ = fmt.Fprintf(c.w, "%s %s err: %v %s\n", c.name, c.arg, err, time.Now().Sub(c.ts))
	return err
}
