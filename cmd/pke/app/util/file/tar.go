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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"emperror.dev/errors"
)

func Untar(out io.Writer, r io.Reader) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "unable to open gzip")
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "unable to read next tar item")
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			hdr.Name = Absolutise(string(os.PathSeparator), hdr.Name)
			err = os.MkdirAll(hdr.Name, os.FileMode(hdr.Mode))
			_, _ = fmt.Fprintf(out, "mkdir -p %s err: %v\n", hdr.Name, err)
			if err != nil {
				return errors.Wrapf(err, "unable to create directory: %s, mode: %v", hdr.Name, hdr.FileInfo().Mode())
			}
		case tar.TypeReg:
			hdr.Name = Absolutise(string(os.PathSeparator), hdr.Name)
			_, _ = fmt.Fprintf(out, "write %s ", hdr.Name)
			f, err := os.OpenFile(hdr.Name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, hdr.FileInfo().Mode())
			if err != nil {
				if e, ok := err.(*os.PathError); ok && e.Err == syscall.EEXIST {
					_, _ = fmt.Fprintf(out, "exist, skipping\n")
					continue
				}
				_, _ = fmt.Fprintf(out, "err: %v\n", err)
				return errors.Wrapf(err, "unable to create file: %s, mode: %v", hdr.Name, hdr.FileInfo().Mode())
			}
			n, err := io.Copy(f, tr)
			_ = f.Close()
			if err != nil {
				_, _ = fmt.Fprintf(out, "err: %v\n", err)
				return errors.Wrapf(err, "unable to write file: %s, mode: %v", hdr.Name, hdr.FileInfo().Mode())
			}
			if hdr.Size != n {
				_, _ = fmt.Fprintf(out, "err: %v\n", err)
				return errors.Wrapf(err, "write failure. file: %s, written: %d, expected: %d", hdr.Name, hdr.Size, n)
			}
			_, _ = fmt.Fprintf(out, "ok\n")
		}
	}

	return nil
}

// Absolutise create absolute filepath from relative with given base path.
func Absolutise(base, p string) string {
	p = filepath.Clean(p)
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(base, filepath.Dir(p), filepath.Base(p))
}
