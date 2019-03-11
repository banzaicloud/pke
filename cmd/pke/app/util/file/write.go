package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func Overwrite(file, contents string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return errors.Wrapf(err, "unable to create %q file", file)
	}
	defer func() { _ = f.Close() }()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrapf(err, "unable to read %q file", file)
	}
	if !bytes.Equal([]byte(contents), b) {
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return errors.Wrapf(err, "unable to seek %q file", file)
		}
		_, err = f.WriteString(contents)
		if err != nil {
			return errors.Wrapf(err, "unable to write %q file", file)
		}

		n, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return errors.Wrapf(err, "unable to seek %q file", file)
		}
		err = f.Truncate(n)
		if err != nil {
			return errors.Wrapf(err, "unable to truncate %q file", file)
		}
	}

	return nil
}
