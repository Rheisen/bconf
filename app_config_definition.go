package bconf

import "github.com/rheisen/bconf/bconfenum"

type AppConfigDefinition struct {
	Name         string
	ConfigFields map[string]Field
	KeyPrefix    string
	Loaders      []bconfenum.ConfigLoader
}

func (d *AppConfigDefinition) Validate() []error {
	errs := []error{}

	if d.ConfigFields != nil {
		for _, field := range d.ConfigFields {
			errs = append(errs, field.Validate()...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (d *AppConfigDefinition) clone() AppConfigDefinition {
	clone := AppConfigDefinition{
		Name:      d.Name,
		KeyPrefix: d.KeyPrefix,
	}

	if len(d.Loaders) > 0 {
		clone.Loaders = make([]bconfenum.ConfigLoader, len(d.Loaders))
		for index, value := range d.Loaders {
			clone.Loaders[index] = value
		}
	}

	if len(d.ConfigFields) > 0 {
		clone.ConfigFields = make(map[string]Field, len(d.ConfigFields))
		for key, value := range d.ConfigFields {
			clone.ConfigFields[key] = value
		}
	}

	return clone
}

func (d *AppConfigDefinition) setDefaults() {
	if d.Name == "" {
		d.Name = "app_config"
	}
	if d.Loaders == nil || len(d.Loaders) == 0 {
		d.Loaders = []bconfenum.ConfigLoader{
			bconfenum.ConfigLoaderEnvironment,
		}
	}
}
