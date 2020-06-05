// Copyright © 2020 Banzai Cloud
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

package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

const configPath = "/etc/banzaicloud/pke.yaml"

// Load reads the configuration from the filesystem or returns the default config if it cannot be found.
func Load() (config Config, err error) {
	if !fileExists(configPath) {
		return Default(), nil
	}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("config read: %w", err)
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, fmt.Errorf("config read: %w", err)
	}

	config = loadDefaults(config)

	return config, nil
}

func loadDefaults(config Config) Config {
	def := Default()

	if config.Kubernetes.Version == "" {
		config.Kubernetes.Version = def.Kubernetes.Version
	}

	if config.ContainerRuntime.Type == "" {
		config.ContainerRuntime.Type = def.ContainerRuntime.Type
	}

	return config
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
