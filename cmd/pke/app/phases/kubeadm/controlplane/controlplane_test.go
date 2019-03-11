package controlplane

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/stretchr/testify/require"
)

func TestWriteKubeadmConfig(t *testing.T) {
	t.SkipNow()
	filename := os.TempDir() + "kubeadm.conf"
	t.Log(filename)
	err := WriteKubeadmConfig(
		os.Stdout,
		filename,
		"192.168.64.11:6443",
		"192.168.64.11:6443",
		"my-cluster",
		"",
		"1.12.2",
		"10.32.0.0/24",
		"10.200.0.0/16",
		constants.CloudProviderAmazon,
		"pool1",
		"/etc/kubernetes/pki/cm-signing-ca.crt",
		[]string{"almafa", "vadkorte"},
		"",
		"",
		"",
	)
	require.NoError(t, err)
	defer func() { _ = os.Remove(filename) }()

	b, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	t.Logf("%s\n", b)
}

func TestWriteKubeadmAmazonConfig(t *testing.T) {
	t.SkipNow()
	filename := os.TempDir() + "aws.conf"
	t.Log(filename)
	err := writeKubeadmAmazonConfig(os.Stdout, filename, constants.CloudProviderAmazon)
	require.NoError(t, err)
	defer func() { _ = os.Remove(filename) }()

	b, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	t.Logf("%s\n", b)
}
