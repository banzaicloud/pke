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
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type Command struct {
	name string
	arg  []string
	w    io.Writer
	ts   time.Time
	*exec.Cmd
}

func Cmd(w io.Writer, name string, arg ...string) *Command {
	return &Command{
		name: name,
		arg:  arg,
		w:    w,
		Cmd:  exec.Command(name, arg...),
	}
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

	c.ts = time.Now()
	stdout, err := c.Cmd.StdoutPipe()
	if err != nil {
		return lastLine, err
	}
	stderr, err := c.Cmd.StderrPipe()
	if err != nil {
		return lastLine, err
	}
	go func() {
		combined := io.MultiReader(stdout, stderr)
		scanner := bufio.NewScanner(combined)
		for scanner.Scan() {
			m := scanner.Text()
			lastLine = m
			_, _ = fmt.Fprintln(c.w, m)
		}
	}()

	err = c.Start()
	if err != nil {
		return lastLine, err
	}

	err = c.Wait()

	return lastLine, err
}

func (c *Command) Output() ([]byte, error) {
	c.ts = time.Now()
	out, err := c.Cmd.Output()
	_, _ = fmt.Fprintf(c.w, "%s %s err: %v %s\n", c.name, c.arg, err, time.Now().Sub(c.ts))
	if len(out) > 0 {
		_, _ = fmt.Fprintln(c.w, string(out))
	}
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
