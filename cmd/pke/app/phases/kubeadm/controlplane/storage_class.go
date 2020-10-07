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

package controlplane

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

func applyDefaultStorageClass(out io.Writer, disableDefaultStorageClass bool, cloudProvider string, azureStorageAccountType, azureStorageKind string, imageRepository string) error {
	if disableDefaultStorageClass {
		return nil
	}

	var err error
	switch cloudProvider {
	case constants.CloudProviderAmazon:
		err = writeStorageClassAmazon(out, storageClassConfig)
	case constants.CloudProviderAzure:
		err = writeStorageClassAzure(out, storageClassConfig, azureStorageAccountType, azureStorageKind)
	case constants.CloudProviderVsphere:
		err = writeStorageClassVsphere(out, storageClassConfig)
	case constants.CloudProviderExternal:
		// TODO: out-of-tree CSI volume plugins
		return nil
	default:
		err = writeStorageClassLocalPathStorage(out, storageClassConfig, imageRepository)
	}
	if err != nil {
		return err
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", storageClassConfig)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err = cmd.CombinedOutputAsync()
	if err != nil {
		return err
	}

	return nil
}

//go:generate templify -t ${GOTMPL} -p controlplane -f storageClassAmazon storage_class_amazon.yaml.tmpl

func writeStorageClassAmazon(out io.Writer, filename string) error {
	_, _ = fmt.Fprintf(out, "[%s] creating Amazon default storage class\n", use)
	// https://kubernetes.io/docs/concepts/storage/storage-classes/#aws-ebs
	tmpl, err := template.New("storage-class-amazon").Parse(storageClassAmazonTemplate())
	if err != nil {
		return err
	}

	type data struct{}

	d := data{}

	return file.WriteTemplate(filename, tmpl, d)
}

func testEnableUUIDTrue(device string) (bool, error) {
	if _, err := os.Stat(device); err != nil {
		return false, err
	}

	if _, err := os.Stat("/usr/lib/udev/scsi_id"); err != nil {
		return false, err
	}

	cmd := exec.Command("/usr/lib/udev/scsi_id", "-g", "-u", "-d", device)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	if err, ok := err.(*exec.ExitError); ok {
		if err.ExitCode() == 1 {
			return false, nil
		}
	}
	return false, err
}

//go:generate templify -t ${GOTMPL} -p controlplane -f storageClassVsphere storage_class_vsphere.yaml.tmpl

func writeStorageClassVsphere(out io.Writer, filename string) error {
	ok, err := testEnableUUIDTrue("/dev/sda")
	switch {
	case err != nil:
		_, _ = fmt.Fprintf(out, "[%s] could not test for disk.EnableUUID=TRUE setting of the VM: %v\n", use, err)
	case !ok:
		return errors.New("It seems like disk.EnableUUID option of this virtual machine is not enabled. Persistent volume vSphere integration won't work without this setting.")
	default:
		_, _ = fmt.Fprintf(out, "[%s] test for disk.EnableUUID=TRUE passed\n", use)
	}

	_, _ = fmt.Fprintf(out, "[%s] creating vSphere default storage class\n", use)
	// https://vmware.github.io/vsphere-storage-for-kubernetes/documentation/policy-based-mgmt.html#vmfs-and-nfs
	tmpl, err := template.New("vshpere-storage-class").Parse(storageClassVsphereTemplate())
	if err != nil {
		return err
	}

	return file.WriteTemplate(filename, tmpl, nil)
}

//go:generate templify -t ${GOTMPL} -p controlplane -f storageClassAzure storage_class_azure.yaml.tmpl

func writeStorageClassAzure(out io.Writer, filename string, storageAccountType, kind string) error {
	_, _ = fmt.Fprintf(out, "[%s] creating Azure default storage class\n", use)
	// https://kubernetes.io/docs/concepts/storage/storage-classes/#new-azure-disk-storage-class-starting-from-v1-7-2
	tmpl, err := template.New("storage-class-azure").Parse(storageClassAzureTemplate())
	if err != nil {
		return err
	}

	type data struct {
		StorageAccountType string
		Kind               string
	}

	d := data{
		StorageAccountType: storageAccountType,
		Kind:               kind,
	}

	return file.WriteTemplate(filename, tmpl, d)
}

//go:generate templify -t ${GOTMPL} -p controlplane -f storageClassLocalPathStorage storage_class_local_path_storage.yaml.tmpl

func writeStorageClassLocalPathStorage(out io.Writer, filename string, imageRepository string) error {
	_, _ = fmt.Fprintf(out, "[%s] creating local default storage class\n", use)

	tmpl, err := template.New("storage-class-local-path").Parse(storageClassLocalPathStorageTemplate())
	if err != nil {
		return err
	}

	type data struct {
		ImageRepository string
	}

	d := data{
		ImageRepository: imageRepository,
	}

	return file.WriteTemplate(filename, tmpl, d)
}
