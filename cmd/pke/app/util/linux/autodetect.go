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

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
)

func KubernetesPackagesImpl(out io.Writer) (KubernetesPackages, error) {
	ver, err := CentOSVersion(out)
	if err != nil {
		ver, err = RedHatVersion(out)
	}
	if err == nil {
		v, _ := semver.NewVersion(ver)
		c7, _ := semver.NewConstraint("7.x-0")
		c8, _ := semver.NewConstraint("8.x-0")
		if c7.Check(v) {
			return NewYumInstaller(), nil
		}
		if c8.Check(v) {
			return NewDnfInstaller(), nil
		}
	}

	if distro, err := LSBReleaseDistributorID(out); err == nil {
		if distro == "Ubuntu" {
			relNum, err := LSBReleaseReleaseNumber(out)
			if err == nil {
				if relNum == "18.04" || relNum == "20.04" {
					return NewAptInstaller(), nil
				}
			}
		}
		return nil, constants.ErrUnsupportedOS
	}

	return nil, constants.ErrUnsupportedOS
}

func ContainerdPackagesImpl(out io.Writer) (ContainerdPackages, error) {
	ver, err := CentOSVersion(out)
	if err != nil {
		ver, err = RedHatVersion(out)
	}
	if err == nil {
		v, _ := semver.NewVersion(ver)
		c7, _ := semver.NewConstraint("7.x-0")
		c8, _ := semver.NewConstraint("8.x-0")
		if c7.Check(v) {
			return NewYumInstaller(), nil
		}
		if c8.Check(v) {
			return NewDnfInstaller(), nil
		}
	}

	if distro, err := LSBReleaseDistributorID(out); err == nil {
		if distro == "Ubuntu" {
			relNum, err := LSBReleaseReleaseNumber(out)
			if err == nil {
				if relNum == "18.04" || relNum == "20.04" {
					return NewAptInstaller(), nil
				}
			}
		}
		return nil, constants.ErrUnsupportedOS
	}

	return nil, constants.ErrUnsupportedOS
}
