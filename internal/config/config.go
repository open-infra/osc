package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/open-infra/osc/internal/client"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// OscConfig represents Osc configuration dir env var.
const OscConfig = "OSCCONFIG"

var (
	// DefaultOscHome represent Osc home directory.
	DefaultOscHome = filepath.Join(mustOscHome(), ".osc")
	// OscConfigFile represents Osc config file location.
	OscConfigFile = filepath.Join(OscHome(), "config.yml")
	// OscLogs represents Osc log.
	OscLogs = filepath.Join(os.TempDir(), fmt.Sprintf("osc-%s.log", MustOscUser()))
	// OscDumpDir represents a directory where Osc screen dumps will be persisted.
	OscDumpDir = filepath.Join(os.TempDir(), fmt.Sprintf("osc-screens-%s", MustOscUser()))
)

type (
	// KubeSettings exposes kubeconfig context information.
	KubeSettings interface {
		// CurrentContextName returns the name of the current context.
		CurrentContextName() (string, error)

		// CurrentClusterName returns the name of the current cluster.
		CurrentClusterName() (string, error)

		// CurrentNamespace returns the name of the current namespace.
		CurrentNamespaceName() (string, error)

		// ClusterNames() returns all available cluster names.
		ClusterNames() ([]string, error)

		// NamespaceNames returns all available namespace names.
		NamespaceNames(nn []v1.Namespace) []string
	}

	// Config tracks Osc configuration options.
	Config struct {
		Osc      *Osc `yaml:"osc"`
		client   client.Connection
		settings KubeSettings
	}
)

// OscHome returns osc configs home directory.
func OscHome() string {
	if env := os.Getenv(OscConfig); env != "" {
		return env
	}

	return DefaultOscHome
}

// NewConfig creates a new default config.
func NewConfig(ks KubeSettings) *Config {
	return &Config{Osc: NewOsc(), settings: ks}
}

// Refine the configuration based on cli args.
func (c *Config) Refine(flags *genericclioptions.ConfigFlags) error {
	cfg, err := flags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	if isSet(flags.Context) {
		c.Osc.CurrentContext = *flags.Context
	} else {
		c.Osc.CurrentContext = cfg.CurrentContext
	}
	log.Debug().Msgf("Active Context %q", c.Osc.CurrentContext)
	if c.Osc.CurrentContext == "" {
		return errors.New("Invalid kubeconfig context detected")
	}
	context, ok := cfg.Contexts[c.Osc.CurrentContext]
	if !ok {
		return fmt.Errorf("The specified context %q does not exists in kubeconfig", c.Osc.CurrentContext)
	}
	c.Osc.CurrentCluster = context.Cluster
	if len(context.Namespace) != 0 {
		if err := c.SetActiveNamespace(context.Namespace); err != nil {
			return err
		}
	}

	if isSet(flags.ClusterName) {
		c.Osc.CurrentCluster = *flags.ClusterName
	}

	if isSet(flags.Namespace) {
		if err := c.SetActiveNamespace(*flags.Namespace); err != nil {
			return err
		}
	}

	return nil
}

// Reset the context to the new current context/cluster.
// if it does not exist.
func (c *Config) Reset() {
	c.Osc.CurrentContext, c.Osc.CurrentCluster = "", ""
}

// CurrentCluster fetch the configuration activeCluster.
func (c *Config) CurrentCluster() *Cluster {
	if c, ok := c.Osc.Clusters[c.Osc.CurrentCluster]; ok {
		return c
	}
	return nil
}

// ActiveNamespace returns the active namespace in the current cluster.
func (c *Config) ActiveNamespace() string {
	if cl := c.CurrentCluster(); cl != nil {
		if cl.Namespace != nil {
			return cl.Namespace.Active
		}
	}
	return "default"
}

// ValidateFavorites ensure favorite ns are legit.
func (c *Config) ValidateFavorites() {
	cl := c.Osc.ActiveCluster()
	if cl == nil {
		cl = NewCluster()
	}
	cl.Validate(c.client, c.settings)
	cl.Namespace.Validate(c.client, c.settings)
}

// FavNamespaces returns fav namespaces in the current cluster.
func (c *Config) FavNamespaces() []string {
	cl := c.Osc.ActiveCluster()
	if cl == nil {
		return nil
	}
	return c.Osc.ActiveCluster().Namespace.Favorites
}

// SetActiveNamespace set the active namespace in the current cluster.
func (c *Config) SetActiveNamespace(ns string) error {
	if c.Osc.ActiveCluster() != nil {
		return c.Osc.ActiveCluster().Namespace.SetActive(ns, c.settings)
	}
	err := errors.New("no active cluster. unable to set active namespace")
	log.Error().Err(err).Msg("SetActiveNamespace")

	return err
}

// ActiveView returns the active view in the current cluster.
func (c *Config) ActiveView() string {
	if c.Osc.ActiveCluster() == nil {
		return defaultView
	}

	cmd := c.Osc.ActiveCluster().View.Active
	if c.Osc.manualCommand != nil && *c.Osc.manualCommand != "" {
		cmd = *c.Osc.manualCommand
	}

	return cmd
}

// SetActiveView set the currently cluster active view
func (c *Config) SetActiveView(view string) {
	cl := c.Osc.ActiveCluster()
	if cl != nil {
		cl.View.Active = view
	}
}

// GetConnection return an api server connection.
func (c *Config) GetConnection() client.Connection {
	return c.client
}

// SetConnection set an api server connection.
func (c *Config) SetConnection(conn client.Connection) {
	c.client = conn
}

// Load Osc configuration from file
func (c *Config) Load(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	c.Osc = NewOsc()

	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return err
	}
	if cfg.Osc != nil {
		c.Osc = cfg.Osc
	}
	if c.Osc.Logger == nil {
		c.Osc.Logger = NewLogger()
	}
	return nil
}

// Save configuration to disk.
func (c *Config) Save() error {
	c.Validate()

	return c.SaveFile(OscConfigFile)
}

// SaveFile Osc configuration to disk.
func (c *Config) SaveFile(path string) error {
	EnsurePath(path, DefaultDirMod)
	cfg, err := yaml.Marshal(c)
	if err != nil {
		log.Error().Msgf("[Config] Unable to save Osc config file: %v", err)
		return err
	}
	return ioutil.WriteFile(path, cfg, 0644)
}

// Validate the configuration.
func (c *Config) Validate() {
	c.Osc.Validate(c.client, c.settings)
}

// Dump debug...
func (c *Config) Dump(msg string) {
	log.Debug().Msgf("Current Cluster: %s\n", c.Osc.CurrentCluster)
	for k, cl := range c.Osc.Clusters {
		log.Debug().Msgf("Osc cluster: %s -- %s\n", k, cl.Namespace)
	}
}

// ----------------------------------------------------------------------------
// Helpers...

func isSet(s *string) bool {
	return s != nil && len(*s) > 0
}
