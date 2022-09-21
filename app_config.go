package bconf

import (
	"fmt"

	"github.com/rheisen/bconf/bconfenum"
)

type AppConfig struct {
	configDefinition *AppConfigDefinition
}

func (c *AppConfig) initialize() []error {
	errs := []error{}
	if errs := c.configDefinition.Validate(); errs != nil && len(errs) > 0 {
		return errs
	}

	c.configDefinition.setDefaults()

	for _, loader := range c.configDefinition.Loaders {
		switch loader {
		case bconfenum.ConfigLoaderEnvironment:

			// os.LookupEnv()
		case bconfenum.ConfigLoaderFlags:
		default:
			errs = append(errs, fmt.Errorf("unsupported config loader found: '%s'", loader.String()))
		}
	}

	return nil
}

func (c *AppConfig) GetField(key string) (Field, error) {
	return Field{}, nil
}

// GetString(key string) (string, error)
// GetInt(key string) (int, error)
// GetBool(key string) (bool, error)
// GetTime(key string) (time.Time, error)
// GetDuration(key string) (time.Duration, error)
