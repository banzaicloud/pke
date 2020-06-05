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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/upgrade"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/flags"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

//go:generate templify -t ${GOTMPL} -p controlplane -f kubeadmConfigV1Beta1 kubeadm_v1beta1.yaml.tmpl
//go:generate templify -t ${GOTMPL} -p controlplane -f kubeadmConfigV1Beta2 kubeadm_v1beta2.yaml.tmpl

const (
	use   = "kubernetes-controlplane"
	short = "Kubernetes Control Plane upgrade"

	kubeConfig                    = "/etc/kubernetes/admin.conf"
	cmdKubeadm                    = "kubeadm"
	cmdKubectl                    = "kubectl"
	certificateAutoApproverUpdate = "/etc/kubernetes/admission-control/deploy-auto-approver-update.yaml"
)

var _ phases.Runnable = (*ControlPlane)(nil)

type ControlPlane struct {
	config config.Config

	kubernetesVersion                string
	kubernetesAdditionalControlPlane bool
	kubeadmConfigMap                 kubeadmConfigMap
	kubeadmConfigUpgrade             string
}

type kubeadmConfigMap struct {
	APIServer struct {
		CertSANs  []string `yaml:"certSANs,omitempty"`
		ExtraArgs struct {
			AdmissionControlConfigFile  string `yaml:"admission-control-config-file"`
			AuditLogMaxage              string `yaml:"audit-log-maxage"`
			AuditLogMaxbackup           string `yaml:"audit-log-maxbackup"`
			AuditLogMaxsize             string `yaml:"audit-log-maxsize"`
			AuditLogPath                string `yaml:"audit-log-path"`
			AuditPolicyFile             string `yaml:"audit-policy-file,omitempty"`
			AuthorizationMode           string `yaml:"authorization-mode"`
			DisableAdmissionPlugins     string `yaml:"disable-admission-plugins"`
			EnableAdmissionPlugins      string `yaml:"enable-admission-plugins"`
			EncryptionProviderConfig    string `yaml:"encryption-provider-config"`
			KubeletCertificateAuthority string `yaml:"kubelet-certificate-authority"`
			Profiling                   string `yaml:"profiling"`
			ServiceAccountLookup        string `yaml:"service-account-lookup"`
			TLSCipherSuites             string `yaml:"tls-cipher-suites"`
			EtcdPrefix                  string `yaml:"etcd-prefix,omitempty"`
			OIDCIssuerURL               string `yaml:"oidc-issuer-url,omitempty"`
			OIDCClientID                string `yaml:"oidc-client-id,omitempty"`
			OIDCUsernameClaim           string `yaml:"oidc-username-claim,omitempty"`
			OIDCUsernamePrefix          string `yaml:"oidc-username-prefix,omitempty"`
			OIDCGroupsClaim             string `yaml:"oidc-groups-claim,omitempty"`
			CloudProvider               string `yaml:"cloud-provider,omitempty"`
			CloudConfig                 string `yaml:"cloud-config,omitempty"`
		} `yaml:"extraArgs"`
		ExtraVolumes []struct {
			HostPath  string `yaml:"hostPath"`
			MountPath string `yaml:"mountPath"`
			Name      string `yaml:"name"`
			PathType  string `yaml:"pathType"`
			ReadOnly  bool   `yaml:"readOnly,omitempty"`
		} `yaml:"extraVolumes"`
		TimeoutForControlPlane string `yaml:"timeoutForControlPlane"`
	} `yaml:"apiServer"`
	APIVersion           string `yaml:"apiVersion"`
	CertificatesDir      string `yaml:"certificatesDir"`
	ClusterName          string `yaml:"clusterName"`
	ControlPlaneEndpoint string `yaml:"controlPlaneEndpoint,omitempty"`
	ControllerManager    struct {
		ExtraArgs struct {
			ClusterName              string `yaml:"cluster-name"`
			FeatureGates             string `yaml:"feature-gates"`
			Profiling                string `yaml:"profiling"`
			TerminatedPodGcThreshold string `yaml:"terminated-pod-gc-threshold"`
			ClusterSigningCertFile   string `yaml:"cluster-signing-cert-file,omitempty"`
			CloudProvider            string `yaml:"cloud-provider,omitempty"`
			CloudConfig              string `yaml:"cloud-config,omitempty"`
		} `yaml:"extraArgs"`
		ExtraVolumes []struct {
			HostPath  string `yaml:"hostPath"`
			MountPath string `yaml:"mountPath"`
			Name      string `yaml:"name"`
			PathType  string `yaml:"pathType"`
			ReadOnly  bool   `yaml:"readOnly,omitempty"`
		} `yaml:"extraVolumes,omitempty"`
	} `yaml:"controllerManager"`
	DNS struct {
		Type string `yaml:"type"`
	} `yaml:"dns"`
	Etcd struct {
		Local struct {
			DataDir   string `yaml:"dataDir"`
			ExtraArgs struct {
				PeerAutoTLS string `yaml:"peer-auto-tls"`
			} `yaml:"extraArgs"`
		} `yaml:"local,omitempty"`
		External struct {
			Endpoints []struct {
				CAFile   string `yaml:"caFile"`
				CertFile string `yaml:"certFile"`
				KeyFile  string `yaml:"keyFile"`
			} `yaml:"enpoints"`
		} `yaml:"external,omitempty"`
	} `yaml:"etcd"`
	ImageRepository   string `yaml:"imageRepository"`
	Kind              string `yaml:"kind"`
	KubernetesVersion string `yaml:"kubernetesVersion"`
	Networking        struct {
		DNSDomain     string `yaml:"dnsDomain"`
		PodSubnet     string `yaml:"podSubnet"`
		ServiceSubnet string `yaml:"serviceSubnet"`
	} `yaml:"networking"`
	Scheduler struct {
		ExtraArgs struct {
			Profiling string `yaml:"profiling"`
		} `yaml:"extraArgs"`
	} `yaml:"scheduler"`
	UseHyperKubeImage bool `yaml:"useHyperKubeImage,omitempty"`
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&ControlPlane{config: config})
}

func (*ControlPlane) Use() string {
	return use
}

func (*ControlPlane) Short() string {
	return short
}

func (c *ControlPlane) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, c.config.Kubernetes.Version, "Kubernetes version")
	// Additional Control Plane
	flags.Bool(constants.FlagAdditionalControlPlane, false, "Treat node as additional control plane")
}

func (c *ControlPlane) Validate(cmd *cobra.Command) error {
	var err error

	c.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: c.kubernetesVersion,
	}); err != nil {
		return err
	}

	c.kubernetesAdditionalControlPlane, err = cmd.Flags().GetBool(constants.FlagAdditionalControlPlane)

	flags.PrintFlags(cmd.OutOrStdout(), c.Use(), cmd.Flags())

	return err
}

func (c *ControlPlane) Run(out io.Writer) error {
	return upgrade.RunWithSkewCheck(out, use, c.kubernetesVersion, kubeConfig, c.upgradeMinor, c.upgradePatch)
}

func (c *ControlPlane) upgradeMinor(out io.Writer, from, to *semver.Version) error {
	_, _ = fmt.Fprintf(out, "[%s] upgrading node from %s to %s\n", use, from, to)

	return c.upgradePatch(out, from, to)
}

func (c *ControlPlane) upgradePatch(out io.Writer, from, to *semver.Version) error {
	_, _ = fmt.Fprintf(out, "[%s] patching node from %s to %s\n", use, from, to)

	return c.upgrade(out, from, to)
}

func (c *ControlPlane) upgrade(out io.Writer, from, to *semver.Version) error {
	pm, err := linux.KubernetesPackagesImpl(out)
	if err != nil {
		return err
	}
	err = pm.InstallKubeadmPackage(out, to.String())
	if err != nil {
		return errors.Wrapf(err, "failed to upgrade kubeadm to version %s", to)
	}

	var args []string
	if c.kubernetesAdditionalControlPlane {
		args = []string{
			"upgrade",
			"node",
		}
		version, _ := semver.NewConstraint("<1.15")
		if version.Check(to) {
			args = append(args, "experimental-control-plane")
		} else {
			args = append(args, "--kubelet-version")
		}

	} else {
		version, _ := semver.NewConstraint(">1.14")
		if version.Check(to) {
			err := c.getKubeadmConfigmap(out)
			if err != nil {
				return err
			}

			err = c.generateNewKubeadmConfig(out, from, to)
			if err != nil {
				return err
			}
			err = c.uploadKubeadmConf(out)
			if err != nil {
				return err
			}
		}

		args = []string{
			"upgrade",
			"apply",
			"-f",
		}
		version, _ = semver.NewConstraint("1.16.x")
		if version.Check(to) {
			args = append(args, "--ignore-preflight-errors=CoreDNSUnsupportedPlugins")
		}
	}
	// target version
	args = append(args, to.String())

	_, err = runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	err = pm.InstallKubernetesPackages(out, to.String())
	if err != nil {
		return err
	}

	err = linux.SystemctlReload(out)
	if err != nil {
		return err
	}

	err = linux.SystemctlStop(out, "kubelet")
	if err != nil {
		return err
	}

	err = linux.SystemctlStart(out, "kubelet")
	if err != nil {
		return err
	}

	fromVersion, _ := semver.NewConstraint("1.17.x")
	toVersion, _ := semver.NewConstraint("1.18.x")
	if fromVersion.Check(from) && toVersion.Check(to) {
		// apply AutoApproverRbacUpdate
		if err := writeCertificateAutoApproverRbacUpdate(out); err != nil {
			return err
		}
	}

	return nil
}

func (c *ControlPlane) uploadKubeadmConf(out io.Writer) error {
	args := []string{
		"init",
		"phase",
		"upload-config",
		"kubeadm",
		"--config",
		c.kubeadmConfigUpgrade,
	}

	_, err := runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	return nil
}

func (c *ControlPlane) getKubeadmConfigmap(out io.Writer) error {
	cmd := runner.Cmd(ioutil.Discard, cmdKubeadm, "config", "view")
	o, err := cmd.Output()
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(o, &c.kubeadmConfigMap)
	if err != nil {
		return err
	}
	return nil
}

func (c *ControlPlane) generateNewKubeadmConfig(out io.Writer, from, to *semver.Version) error {
	var conf string
	switch to.Minor() {
	case 15, 16, 17:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1
		c.kubeadmConfigMap.UseHyperKubeImage = true
		conf = kubeadmConfigV1Beta1Template()
	case 18:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2
		conf = kubeadmConfigV1Beta2Template()
	default:
		return errors.New(fmt.Sprintf("unsupported Kubernetes version %q for upgrade", c.kubernetesVersion))
	}

	version, _ := semver.NewConstraint(">1.17")
	if version.Check(to) {
		c.kubeadmConfigMap.UseHyperKubeImage = false
		admissionPlugins := strings.Split(c.kubeadmConfigMap.APIServer.ExtraArgs.EnableAdmissionPlugins, ",")
		upgradeAddmissionPlugins := []string{}
		for _, plugin := range admissionPlugins {
			if plugin != "DenyEscalatingExec" {
				upgradeAddmissionPlugins = append(upgradeAddmissionPlugins, plugin)
			}
		}
		c.kubeadmConfigMap.APIServer.ExtraArgs.EnableAdmissionPlugins = strings.Join(upgradeAddmissionPlugins, ",")
	}

	type data struct {
		KubeadmConfig             kubeadmConfigMap
		APIServerAdvertiseAddress string
		APIServerBindPort         string
	}

	d := data{
		KubeadmConfig: c.kubeadmConfigMap,
	}
	apiServerAddress := []string{}
	if c.kubeadmConfigMap.ControlPlaneEndpoint != "" {
		apiServerAddress = strings.Split(c.kubeadmConfigMap.ControlPlaneEndpoint, ":")
		d.APIServerAdvertiseAddress = apiServerAddress[0]
		d.APIServerBindPort = apiServerAddress[1]
	}

	tmpl, err := template.New("kubeadm-config-ugrade").Parse(conf)
	if err != nil {
		return err
	}
	c.kubeadmConfigUpgrade = "/etc/kubernetes/kubeadm-" + time.Now().Format(time.RFC3339) + ".conf"

	return file.WriteTemplate(c.kubeadmConfigUpgrade, tmpl, d)
}

//go:generate templify -t ${GOTMPL} -p controlplane -f certificateAutoApproverRbacUpdate certificate_auto_approver_rbac_update.yaml.tmpl

func writeCertificateAutoApproverRbacUpdate(out io.Writer) error {
	filename := certificateAutoApproverUpdate
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	err = file.Overwrite(filename, certificateAutoApproverRbacUpdateTemplate())
	if err != nil {
		return err
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", filename)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err = cmd.CombinedOutputAsync()
	return err
}
