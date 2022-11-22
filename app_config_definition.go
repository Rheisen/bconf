package bconf

import (
	"fmt"
	"os"
	"strings"

	"github.com/rheisen/bconf/bconfconst"
)

type AppConfigDefinition struct {
	Name         string
	ConfigFields map[string]*Field
	KeyPrefix    string
	Loaders      []string // in order
}

func (d *AppConfigDefinition) GenerateFieldDefaults() []error {
	errs := []error{}

	if d.ConfigFields != nil {
		for _, field := range d.ConfigFields {
			if err := field.GenerateDefault(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func (d *AppConfigDefinition) Validate() []error {
	errs := []error{}

	if d.Loaders != nil {
		for _, loader := range d.Loaders {
			if _, found := bconfconst.LoadersMap()[loader]; !found {
				errs = append(errs, fmt.Errorf(
					"invalid loader, expected one-of '%v', found '%s'",
					bconfconst.Loaders(),
					loader,
				))
			}
		}
	}

	if d.ConfigFields != nil {
		for _, field := range d.ConfigFields {
			if err := field.GenerateDefault(); err != nil {
				errs = append(errs, err)
			}
			errs = append(errs, field.Validate()...)
		}
	}

	return errs
}

func (d *AppConfigDefinition) clone() AppConfigDefinition {
	clone := AppConfigDefinition{
		Name:      d.Name,
		KeyPrefix: d.KeyPrefix,
	}

	if len(d.Loaders) > 0 {
		clone.Loaders = make([]string, len(d.Loaders))
		for index, value := range d.Loaders {
			clone.Loaders[index] = value
		}
	}

	if len(d.ConfigFields) > 0 {
		clone.ConfigFields = make(map[string]*Field, len(d.ConfigFields))
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
		d.Loaders = []string{
			bconfconst.EnvironmentLoader,
		}
	}
}

func (d *AppConfigDefinition) loadFields() []error {
	errs := []error{}

	for _, loader := range d.Loaders {
		switch loader {
		case bconfconst.EnvironmentLoader:
			for key, field := range d.ConfigFields {
				envKey := ""
				if d.KeyPrefix != "" {
					envKey = fmt.Sprintf("%s_%s", d.KeyPrefix, key)
				} else {
					envKey = key
				}

				envKey = strings.ToUpper(envKey)
				value, found := os.LookupEnv(envKey)
				if found {
					if err := field.set(bconfconst.EnvironmentLoader, value); err != nil {
						errs = append(errs, fmt.Errorf("invalid field value: %w", err))
					}
				}
			}
		case bconfconst.FlagLoader:
		default:
			errs = append(errs, fmt.Errorf("unsupported loader, found '%s'", loader))
		}
	}

	return errs
}
