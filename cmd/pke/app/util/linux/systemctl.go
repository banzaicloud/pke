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

package linux

import (
	"io"
	"os/exec"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdSystemctl = "/bin/systemctl"
	start        = "start"
	stop         = "stop"
	enable       = "enable"
	disable      = "disable"
	isEnabled    = "is-enabled"
	isActive     = "is-active"
	reload       = "daemon-reload"
)

func Systemctl(out io.Writer, command, service string) error {
	if service != "" {
		return runner.Cmd(out, cmdSystemctl, command, service).Run()
	}

	return runner.Cmd(out, cmdSystemctl, command).Run()
}

func SystemctlReload(out io.Writer) error {
	return Systemctl(out, reload, "")
}

func SystemctlEnable(out io.Writer, service string) error {
	if err := SystemctlReload(out); err != nil {
		return err
	}
	return Systemctl(out, enable, service)
}

func SystemctlDisable(out io.Writer, service string) error {
	if err := SystemctlReload(out); err != nil {
		return err
	}
	return Systemctl(out, disable, service)
}

func SystemctlEnabled(out io.Writer, service string) (bool, error) {
	err := Systemctl(out, isEnabled, service)
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func SystemctlStart(out io.Writer, service string) error {
	if err := SystemctlReload(out); err != nil {
		return err
	}
	return Systemctl(out, start, service)
}

func SystemctlStop(out io.Writer, service string) error {
	if err := SystemctlReload(out); err != nil {
		return err
	}
	return Systemctl(out, stop, service)
}

func SystemctlActive(out io.Writer, service string) (bool, error) {
	err := Systemctl(out, isActive, service)
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func SystemctlEnableAndStart(out io.Writer, service string) error {
	if err := SystemctlEnable(out, service); err != nil {
		return err
	}

	return SystemctlStart(out, service)
}

func SystemctlDisableAndStop(out io.Writer, service string) error {
	if err := SystemctlDisable(out, service); err != nil {
		return err
	}

	return SystemctlStop(out, service)
}
