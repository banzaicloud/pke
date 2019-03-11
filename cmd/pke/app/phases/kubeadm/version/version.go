package version

import (
	"fmt"
	"io"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-version"
	short = "Check Kubernetes version is supported or not"

	constraint = ">= 1.12.0, < 1.12.6 || >= 1.13.0, <= 1.13.4"
)

var _ phases.Runnable = (*Version)(nil)

type Version struct {
	kubernetesVersion string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &Version{})
}

func (v *Version) Use() string {
	return use
}

func (v *Version) Short() string {
	return short
}

func (v *Version) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, "1.13.3", "Kubernetes version")
}

func (v *Version) Validate(cmd *cobra.Command) error {
	var err error
	v.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: v.kubernetesVersion,
	}); err != nil {
		return err
	}

	return validVersion(v.kubernetesVersion, constraint)
}

func (v *Version) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "Kubernetes version %q is supported\n", v.kubernetesVersion)
	return nil
}

func validVersion(version, constraint string) error {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return err
	}

	ver, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	if !c.Check(ver) {
		return errors.Wrapf(constants.ErrUnsupportedKubernetesVersion, "got: %q, expected: %q", version, constraint)
	}

	return nil
}
