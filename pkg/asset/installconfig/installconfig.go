package installconfig

import (
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/conversion"
	"github.com/openshift/installer/pkg/types/defaults"
	openstackvalidation "github.com/openshift/installer/pkg/types/openstack/validation"
	"github.com/openshift/installer/pkg/types/validation"
)

const (
	installConfigFilename = "install-config.yaml"
)

// InstallConfig generates the install-config.yaml file.
type InstallConfig struct {
	Config *types.InstallConfig `json:"config"`
	File   *asset.File          `json:"file"`
}

var _ asset.WritableAsset = (*InstallConfig)(nil)

// Dependencies returns all of the dependencies directly needed by an
// InstallConfig asset.
func (a *InstallConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		&sshPublicKey{},
		&baseDomain{},
		&clusterName{},
		&pullSecret{},
		&platform{},
	}
}

// Generate generates the install-config.yaml file.
func (a *InstallConfig) Generate(parents asset.Parents) error {
	sshPublicKey := &sshPublicKey{}
	baseDomain := &baseDomain{}
	clusterName := &clusterName{}
	pullSecret := &pullSecret{}
	platform := &platform{}
	parents.Get(
		sshPublicKey,
		baseDomain,
		clusterName,
		pullSecret,
		platform,
	)

	a.Config = &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName.ClusterName,
		},
		SSHKey:     sshPublicKey.Key,
		BaseDomain: baseDomain.BaseDomain,
		PullSecret: pullSecret.PullSecret,
	}

	a.Config.AWS = platform.AWS
	a.Config.Libvirt = platform.Libvirt
	a.Config.None = platform.None
	a.Config.OpenStack = platform.OpenStack
	a.Config.VSphere = platform.VSphere
	a.Config.Azure = platform.Azure
	a.Config.GCP = platform.GCP
	a.Config.BareMetal = platform.BareMetal

	return a.finish("")
}

// Name returns the human-friendly name of the asset.
func (a *InstallConfig) Name() string {
	return "Install Config"
}

// Files returns the files generated by the asset.
func (a *InstallConfig) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// Load returns the installconfig from disk.
func (a *InstallConfig) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(installConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &types.InstallConfig{}
	if err := yaml.Unmarshal(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", installConfigFilename)
	}
	a.Config = config

	// Upconvert any deprecated fields
	if err := conversion.ConvertInstallConfig(a.Config); err != nil {
		return false, errors.Wrap(err, "failed to upconvert install config")
	}

	err = a.finish(installConfigFilename)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *InstallConfig) finish(filename string) error {
	defaults.SetInstallConfigDefaults(a.Config)

	if err := validation.ValidateInstallConfig(a.Config, openstackvalidation.NewValidValuesFetcher()).ToAggregate(); err != nil {
		if filename == "" {
			return errors.Wrap(err, "invalid install config")
		}
		return errors.Wrapf(err, "invalid %q file", filename)
	}

	data, err := yaml.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal InstallConfig")
	}
	a.File = &asset.File{
		Filename: installConfigFilename,
		Data:     data,
	}
	return nil
}
