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

package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdKubectl = "/bin/kubectl"
)

type Token struct {
	Token   string    `json:"token"`
	TTL     int       `json:"-"`
	Expires time.Time `json:"expires"`
	Expired bool      `json:"expired"`
}

type Output struct {
	Tokens []*Token `json:"tokens"`
}

func Get(out io.Writer, secret string) (*Token, error) {
	// kubectl get -n kube-system secret -o json
	subCmd := runner.Cmd(out, cmdKubectl, []string{"get", "secret", "-n", "kube-system", "-o", "json", secret}...)
	cmdOut, err := subCmd.Output()
	if err != nil {
		return nil, err
	}
	s := struct {
		Data map[string]string `json:"data"`
	}{}
	err = json.Unmarshal(cmdOut, &s)
	if err != nil {
		return nil, err
	}

	tid, err := base64.StdEncoding.DecodeString(s.Data["token-id"])
	if err != nil {
		return nil, err
	}
	ts, err := base64.StdEncoding.DecodeString(s.Data["token-secret"])
	if err != nil {
		return nil, err
	}
	exp, err := base64.StdEncoding.DecodeString(s.Data["expiration"])
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, string(exp))
	if err != nil {
		return nil, err
	}

	return &Token{
		Token:   fmt.Sprintf("%s.%s", tid, ts),
		TTL:     int(t.Sub(time.Now()).Hours()),
		Expires: t,
		Expired: t.Sub(time.Now()) <= 0,
	}, nil
}
