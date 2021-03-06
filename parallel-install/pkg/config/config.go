//Package config defines top-level configuration settings for library users.
package config

import (
	"fmt"
	"os"
	"time"
)

//Configures various install/uninstall operation parameters.
//There are no different parameters for the "install" and "delete" operations.
//If you need different configurations, just use two different Installation instances.
type Config struct {
	//Number of parallel workers used for an install/uninstall operation
	WorkersCount int
	//After this time workers' context is canceled. Pending worker goroutines (if any) may continue if blocked by Helm client.
	CancelTimeout time.Duration
	//After this time install/delete operation is aborted and returns an error to the user.
	//Worker goroutines may still be working in the background.
	//Must be greater than CancelTimeout.
	QuitTimeout time.Duration
	//Timeout for the underlying Helm client
	HelmTimeoutSeconds int
	//Initial interval used for exponent backoff retry policy
	BackoffInitialIntervalSeconds int
	//Maximum time used for exponent backoff retry policy
	BackoffMaxElapsedTimeSeconds int
	//Logger to use
	Log func(format string, v ...interface{})
	//Maximum number of Helm revision saved per release
	HelmMaxRevisionHistory int
	//Installation / Upgrade profile: evaluation|production
	Profile string
	// Path to Kyma components list
	ComponentsListFile string
	// Path to Kyma resources
	ResourcePath string
	// Path to Kyma CRDs
	CrdPath string
	//Kyma version
	Version string
}

// Validate the given configuration options
func (c *Config) Validate() error {
	if c.WorkersCount <= 0 {
		return fmt.Errorf("Workers count cannot be <= 0")
	}
	if err := c.pathExists(c.ComponentsListFile, "Components list"); err != nil {
		return err
	}
	if err := c.pathExists(c.ResourcePath, "Resource path"); err != nil {
		return err
	}
	if err := c.pathExists(c.CrdPath, "CRD path"); err != nil {
		return err
	}
	if c.Version == "" {
		return fmt.Errorf("Version is empty")
	}
	return nil
}

func (c *Config) pathExists(path string, description string) error {
	if path == "" {
		return fmt.Errorf("%s is empty", description)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s '%s' not found", description, path)
	}
	return nil
}
