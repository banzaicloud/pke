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
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdKubectl = "kubectl"
	kubeConfig = "/etc/kubernetes/admin.conf"
)

type Token struct {
	Token    string    `json:"token"`
	TTL      int       `json:"-"`
	Expires  time.Time `json:"expires"`
	Expired  bool      `json:"expired"`
	CertHash string    `json:"hash"`
}

type Output struct {
	Tokens []*Token `json:"tokens"`
}

func Get(out io.Writer, secret, certHash string) (*Token, error) {
	// kubectl get -n kube-system secret -o json
	subCmd := runner.Cmd(out, cmdKubectl, []string{"get", "secret", "-n", "kube-system", "-o", "json", secret}...)
	subCmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
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

	var (
		t       time.Time
		ttl     int
		expired bool
	)
	if exps := string(exp); exps != "" {
		t, err = time.Parse(time.RFC3339, exps)
		if err != nil {
			return nil, err
		}
		ttl = int(t.Sub(time.Now()).Hours())
		expired = t.Sub(time.Now()) <= 0
	}

	return &Token{
		Token:    fmt.Sprintf("%s.%s", tid, ts),
		TTL:      ttl,
		Expires:  t,
		Expired:  expired,
		CertHash: certHash,
	}, nil
}

func CertHash(out io.Writer, certFile string) (string, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to open certificate for hashing")
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return "", errors.New("failed to parse certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse certificate")
	}

	h := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(h[:])), nil
}
