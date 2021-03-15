package config

import (
	"github.com/open-infra/osc/internal/client"
)

const (
	defaultRefreshRate  = 2
	defaultMaxConnRetry = 5
)

// Osc tracks Osc configuration options.
type Osc struct {
	RefreshRate       int                 `yaml:"refreshRate"`
	MaxConnRetry      int                 `yaml:"maxConnRetry"`
	EnableMouse       bool                `yaml:"enableMouse"`
	Headless          bool                `yaml:"headless"`
	Crumbsless        bool                `yaml:"crumbsless"`
	ReadOnly          bool                `yaml:"readOnly"`
	NoIcons           bool                `yaml:"noIcons"`
	Logger            *Logger             `yaml:"logger"`
	CurrentContext    string              `yaml:"currentContext"`
	CurrentCluster    string              `yaml:"currentCluster"`
	Clusters          map[string]*Cluster `yaml:"clusters,omitempty"`
	Thresholds        Threshold           `yaml:"thresholds"`
	manualRefreshRate int
	manualHeadless    *bool
	manualCrumbsless  *bool
	manualReadOnly    *bool
	manualCommand     *string
}

// NewOsc create a new K9s configuration.
func NewOsc() *Osc {
	return &Osc{
		RefreshRate:  defaultRefreshRate,
		MaxConnRetry: defaultMaxConnRetry,
		Logger:       NewLogger(),
		Clusters:     make(map[string]*Cluster),
		Thresholds:   NewThreshold(),
	}
}

// OverrideRefreshRate set the refresh rate manually.
func (k *Osc) OverrideRefreshRate(r int) {
	k.manualRefreshRate = r
}

// OverrideHeadless set the headlessness manually.
func (k *Osc) OverrideHeadless(b bool) {
	k.manualHeadless = &b
}

// OverrideCrumbsless set the headlessness manually.
func (k *Osc) OverrideCrumbsless(b bool) {
	k.manualCrumbsless = &b
}

// OverrideReadOnly set the readonly mode manually.
func (k *Osc) OverrideReadOnly(b bool) {
	if b {
		k.manualReadOnly = &b
	}
}

// OverrideWrite set the write mode manually.
func (k *Osc) OverrideWrite(b bool) {
	if b {
		var flag bool
		k.manualReadOnly = &flag
	}
}

// OverrideCommand set the command manually.
func (k *Osc) OverrideCommand(cmd string) {
	k.manualCommand = &cmd
}

// IsHeadless returns headless setting.
func (k *Osc) IsHeadless() bool {
	h := k.Headless
	if k.manualHeadless != nil && *k.manualHeadless {
		h = *k.manualHeadless
	}

	return h
}

// IsCrumbsless returns crumbsless setting.
func (k *Osc) IsCrumbsless() bool {
	h := k.Crumbsless
	if k.manualCrumbsless != nil && *k.manualCrumbsless {
		h = *k.manualCrumbsless
	}

	return h
}

// GetRefreshRate returns the current refresh rate.
func (k *Osc) GetRefreshRate() int {
	rate := k.RefreshRate
	if k.manualRefreshRate != 0 {
		rate = k.manualRefreshRate
	}

	return rate
}

// IsReadOnly returns the readonly setting.
func (k *Osc) IsReadOnly() bool {
	readOnly := k.ReadOnly
	if k.manualReadOnly != nil {
		readOnly = *k.manualReadOnly
	}

	return readOnly
}

// ActiveCluster returns the currently active cluster.
func (k *Osc) ActiveCluster() *Cluster {
	if k.Clusters == nil {
		k.Clusters = map[string]*Cluster{}
	}

	if c, ok := k.Clusters[k.CurrentCluster]; ok {
		return c
	}
	k.Clusters[k.CurrentCluster] = NewCluster()

	return k.Clusters[k.CurrentCluster]
}

func (k *Osc) validateDefaults() {
	if k.RefreshRate <= 0 {
		k.RefreshRate = defaultRefreshRate
	}
	if k.MaxConnRetry <= 0 {
		k.MaxConnRetry = defaultMaxConnRetry
	}
}

func (k *Osc) validateClusters(c client.Connection, ks KubeSettings) {
	cc, err := ks.ClusterNames()
	if err != nil {
		return
	}
	for key := range k.Clusters {
		k.Clusters[key].Validate(c, ks)
		if InList(cc, key) {
			continue
		}
		if k.CurrentCluster == key {
			k.CurrentCluster = ""
		}
		delete(k.Clusters, key)
	}
}

// Validate the current configuration.
func (k *Osc) Validate(c client.Connection, ks KubeSettings) {
	k.validateDefaults()
	if k.Clusters == nil {
		k.Clusters = map[string]*Cluster{}
	}
	k.validateClusters(c, ks)

	if k.Logger == nil {
		k.Logger = NewLogger()
	} else {
		k.Logger.Validate(c, ks)
	}
	if k.Thresholds == nil {
		k.Thresholds = NewThreshold()
	}
	k.Thresholds.Validate(c, ks)

	if ctx, err := ks.CurrentContextName(); err == nil && len(k.CurrentContext) == 0 {
		k.CurrentContext = ctx
		k.CurrentCluster = ""
	}

	if cl, err := ks.CurrentClusterName(); err == nil && len(k.CurrentCluster) == 0 {
		k.CurrentCluster = cl
	}

	if _, ok := k.Clusters[k.CurrentCluster]; !ok {
		k.Clusters[k.CurrentCluster] = NewCluster()
	}
	k.Clusters[k.CurrentCluster].Validate(c, ks)
}
