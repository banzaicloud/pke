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

package file

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"
)

func TestDownloadWithSHA256(t *testing.T) {
	// t.SkipNow()
	testCases := []struct {
		u    string
		hash string
		len  int
	}{
		{
			"https://storage.googleapis.com/cri-containerd-release/cri-containerd-1.2.0.linux-amd64.tar.gz",
			"ee076c6260de140f9aa6dee30b0e360abfb80af252d271e697982d1209ca5dee",
			45449776,
		},
	}

	for _, tc := range testCases {
		u, err := url.Parse(tc.u)
		if err != nil {
			t.Error(err)
		}
		f, err := ioutil.TempFile("", "download_test")
		if err != nil {
			t.Error(err)
		}
		defer func() { _ = os.Remove(f.Name()) }()

		err = Download(u, f.Name())
		_ = f.Close()
		if err != nil {
			t.Error(err)
		}

		err = SHA256File(f.Name(), tc.hash)
		if err != nil {
			t.Error(err)
		}

		fi, err := os.Stat(f.Name())
		if err != nil {
			t.Error(err)
		}
		if fi.Size() != int64(tc.len) {
			t.Errorf("unexpected length. got: %d, expected %d", fi.Size(), tc.len)
		}
	}
}
