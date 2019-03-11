package file

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

func Download(u *url.URL, f string) error {
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unhandled http status code: %d", resp.StatusCode)
	}
	defer func() { _ = resp.Body.Close() }()

	out, err := os.Create(f)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func SHA256(filepath string) (string, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err = h.Write(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func SHA256File(f, hash string) error {
	hs, err := SHA256(f)
	if err != nil {
		return err
	}
	if hs != hash {
		return errors.Errorf("hash mismatch. got: %q, expected: %q", hs, hash)
	}

	return nil
}
