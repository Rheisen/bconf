package bconf

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type AppConfigDefinition struct {
	// Name is used in --help text output to identify the application.
	Name string
	// Description is used in --help text output to describe the application.
	Description string
	// ConfigFields defines a mapping to configuration fields.
	ConfigFields map[string]*Field
	// Loaders defines where field values are extracted from. Values found from later loaders take presedence.
	Loaders []Loader
	// HandleHelpFlag defines if the app configuration should handle the --help flag.
	HandleHelpFlag bool
}

func (d *AppConfigDefinition) Clone() AppConfigDefinition {
	clone := *d

	if len(d.Loaders) > 0 {
		clone.Loaders = make([]Loader, len(d.Loaders))
		for index, value := range d.Loaders {
			clone.Loaders[index] = value.Clone()
		}
	}

	if len(d.ConfigFields) > 0 {
		clone.ConfigFields = make(map[string]*Field, len(d.ConfigFields))
		for key, field := range d.ConfigFields {
			newField := *field
			clone.ConfigFields[key] = &newField
		}
	}

	return clone
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
		loaderNames := make(map[string]struct{}, len(d.Loaders))
		for _, loader := range d.Loaders {
			if _, found := loaderNames[loader.Name()]; found {
				errs = append(errs, fmt.Errorf("duplicate loader name found: '%s'", loader.Name()))
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

func (d *AppConfigDefinition) HelpText() string {
	builder := strings.Builder{}

	if d.Name != "" {
		builder.WriteString(fmt.Sprintf("\nUsage of '%s':\n", d.Name))
	} else {
		builder.WriteString(fmt.Sprintf("\nUsage of '%s':\n", os.Args[0]))
	}

	if d.Description != "" {
		builder.WriteString(fmt.Sprintf("%s\n\n", d.Description))
	}

	if len(d.ConfigFields) > 0 {
		keys := make([]string, len(d.ConfigFields))
		idx := 0
		for key := range d.ConfigFields {
			keys[idx] = key
			idx = idx + 1
		}
		sort.Strings(keys)

		requiredFields := []string{}
		optionalFields := []string{}
		for _, key := range keys {
			field := d.ConfigFields[key]
			if field.Required {
				requiredFields = append(requiredFields, key)
			} else {
				optionalFields = append(optionalFields, key)
			}
		}

		if len(requiredFields) > 0 {
			builder.WriteString("Required Configuration:\n")
			for _, key := range requiredFields {
				builder.WriteString(fmt.Sprintf("\t%s", d.fieldHelpString(key)))
			}
		}

		if len(optionalFields) > 0 {
			builder.WriteString("Optional Configuration:\n")
			for _, key := range optionalFields {
				builder.WriteString(fmt.Sprintf("\t%s", d.fieldHelpString(key)))
			}
		}
	}

	return builder.String()
}

func (d *AppConfigDefinition) loadFields() []error {
	errs := []error{}

	for _, loader := range d.Loaders {
		for key, field := range d.ConfigFields {
			value, found := loader.Get(key)
			if found {
				if err := field.set(loader.Name(), value); err != nil {
					errs = append(errs, fmt.Errorf("invalid field value: %w", err))
				}
			}
		}
	}

	return errs
}

func (d *AppConfigDefinition) printHelpText() {
	fmt.Printf("%s", d.HelpText())
}

func (d *AppConfigDefinition) fieldHelpString(key string) string {
	field := d.ConfigFields[key]
	if field == nil {
		return ""
	}

	builder := strings.Builder{}
	spaceBuffer := "\t\t"

	builder.WriteString(fmt.Sprintf("%s %s\n", key, field.FieldType))
	if field.Description != "" {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("%s\n", field.Description))
	}
	if len(field.Enumeration) > 0 {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Accepted values: %s\n", field.enumerationString()))
	}
	if field.Default != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Default value: '%v'\n", field.Default))
	}
	if field.DefaultGenerator != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString("Default value: <generated-at-run-time>\n")
	}

	for _, loader := range d.Loaders {
		helpString := loader.HelpString(key)
		if helpString != "" {
			builder.WriteString(spaceBuffer)
			builder.WriteString(fmt.Sprintf("%s\n", helpString))
		}
	}

	return builder.String()
}
