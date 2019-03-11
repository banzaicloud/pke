package file

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"
)

func TestDownloadWithSHA256(t *testing.T) {
	//t.SkipNow()
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
		if err != nil {
			t.Error(err)
		}

		err = SHA256File(f.Name(), tc.hash)
		if err != nil {
			t.Error(err)
		}

		b, err := ioutil.ReadFile(f.Name())
		if len(b) != tc.len {
			t.Errorf("unexpected length. got: %d, expected %d", len(b), tc.len)
		}
	}
}
