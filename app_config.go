package bconf

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

func NewAppConfig(appName string, appDescription string) *AppConfig {
	return &AppConfig{
		appName:        appName,
		appDescription: appDescription,
		fieldSets:      map[string]*FieldSet{},
		loaders:        []Loader{},
	}
}

type AppConfig struct {
	appName        string
	appDescription string
	fieldSets      map[string]*FieldSet
	fieldSetLock   sync.Mutex
	loaders        []Loader
	register       sync.Once
	registered     bool
}

func (c *AppConfig) AppName() string {
	return c.appName
}

func (c *AppConfig) AppDescription() string {
	return c.appDescription
}

func (c *AppConfig) SetLoaders(loaders ...Loader) []error {
	errs := []error{}

	clonedLoaders := make([]Loader, len(loaders))
	for index, loader := range loaders {
		clonedLoaders[index] = loader.Clone()
	}

	loaderNames := make(map[string]struct{}, len(clonedLoaders))
	for _, loader := range clonedLoaders {
		if _, found := loaderNames[loader.Name()]; found {
			errs = append(errs, fmt.Errorf("duplicate loader name found: '%s'", loader.Name()))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	c.loaders = clonedLoaders

	return nil
}

func (c *AppConfig) AddFieldSet(fieldSet *FieldSet) []error {
	c.fieldSetLock.Lock()
	defer c.fieldSetLock.Unlock()

	errs := []error{}
	fieldSet = fieldSet.Clone()

	// check for field set structural integrity
	if fieldSetErrs := fieldSet.validate(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(errs, fmt.Errorf("field-set '%s' validation error: %w", fieldSet.Key, err))
		}
		return errs
	}

	fieldSet.initializeFieldMap()

	// generate field-set field default values
	if fieldSetErrs := fieldSet.generateFieldDefaults(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field default value generation error: %w", fieldSet.Key, err),
			)
		}
		return errs
	}

	// validate field-set fields
	if fieldSetErrs := fieldSet.validateFields(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field validation error: %w", fieldSet.Key, err),
			)
		}
		return errs
	}

	// persist the field-set to AppConfig
	if c.fieldSets == nil {
		c.fieldSets = map[string]*FieldSet{fieldSet.Key: fieldSet}
		return nil
	}

	if _, keyFound := c.fieldSets[fieldSet.Key]; keyFound {
		errs = append(
			errs,
			fmt.Errorf("duplicate field-set key found: '%s'", fieldSet.Key),
		)
		return errs
	}

	c.fieldSets[fieldSet.Key] = fieldSet

	return nil
}

func (c *AppConfig) AddField(fieldSetKey string, field *Field) []error {
	return nil
}

func (c *AppConfig) LoadFieldSet(fieldSetKey string) []error {
	errs := []error{}

	if !c.registered {
		errs = append(errs, fmt.Errorf("LoadFieldSet cannot be called before the app-config has been registered"))
	}

	if _, fieldSetFound := c.fieldSets[fieldSetKey]; !fieldSetFound {
		errs = append(errs, fmt.Errorf("field-set with key '%s' not found", fieldSetKey))
		return errs
	}

	for _, loader := range c.loaders {
		for key, field := range c.fieldSets[fieldSetKey].fieldMap {
			value, found := loader.Get(fmt.Sprintf("%s_%s", fieldSetKey, key))
			if found {
				if err := field.set(loader.Name(), value); err != nil {
					errs = append(errs, fmt.Errorf("field '%s' load error: %w", key, err))
				}
			}
		}
	}

	return errs
}

func (c *AppConfig) LoadField(fieldSetKey, fieldKey string) []error {
	errs := []error{}

	if !c.registered {
		errs = append(errs, fmt.Errorf("LoadField cannot be called before the app-config has been registered"))
	}

	if _, fieldSetFound := c.fieldSets[fieldSetKey]; !fieldSetFound {
		errs = append(errs, fmt.Errorf("field-set with key '%s' not found", fieldSetKey))
		return errs
	}

	if _, fieldKeyFound := c.fieldSets[fieldSetKey].fieldMap[fieldKey]; !fieldKeyFound {
		errs = append(errs, fmt.Errorf("field with key '%s' not found", fieldKey))
		return errs
	}

	return nil
}

func (c *AppConfig) SetField(fieldSetKey, fieldKey string, fieldValue any) []error {
	return nil
}

// Register loads all defined field sets and optionally checks for and handles the help flag -h and --help.
func (c *AppConfig) Register(handleHelpFlag bool) []error {
	if handleHelpFlag && len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		c.printHelpString()
		os.Exit(0)
	}

	errs := []error{}

	for fieldSetKey := range c.fieldSets {
		if fieldSetErrs := c.LoadFieldSet(fieldSetKey); len(fieldSetErrs) > 0 {
			errs = append(errs, fieldSetErrs...)
		}
	}

	return errs
}

func (c *AppConfig) HelpString() string {
	builder := strings.Builder{}

	if c.appName != "" {
		builder.WriteString(fmt.Sprintf("Usage of '%s':\n", c.appName))
	} else {
		builder.WriteString(fmt.Sprintf("Usage of '%s':\n", os.Args[0]))
	}

	if c.appDescription != "" {
		builder.WriteString(fmt.Sprintf("%s\n\n", c.appDescription))
	}

	fields := c.fields()
	if len(fields) > 0 {
		keys := make([]string, len(fields))
		idx := 0
		for key := range fields {
			keys[idx] = key
			idx += 1
		}
		sort.Strings(keys)

		requiredFields := []string{}
		optionalFields := []string{}
		for _, key := range keys {
			field := fields[key]
			if field.Required {
				requiredFields = append(requiredFields, key)
			} else {
				optionalFields = append(optionalFields, key)
			}
		}

		if len(requiredFields) > 0 {
			builder.WriteString("Required Configuration:\n")
			for _, key := range requiredFields {
				builder.WriteString(fmt.Sprintf("\t%s", c.fieldHelpString(fields, key)))
			}
		}

		if len(optionalFields) > 0 {
			builder.WriteString("Optional Configuration:\n")
			for _, key := range optionalFields {
				builder.WriteString(fmt.Sprintf("\t%s", c.fieldHelpString(fields, key)))
			}
		}
	}

	return builder.String()
}

func (c *AppConfig) fields() map[string]*Field {
	fields := map[string]*Field{}

	for fieldSetKey, fieldSet := range c.fieldSets {
		for _, field := range fieldSet.fieldMap {
			fields[fmt.Sprintf("%s_%s", fieldSetKey, field.Key)] = field
		}
	}

	return fields
}

func (c *AppConfig) fieldHelpString(fields map[string]*Field, key string) string {
	field := fields[key]
	if field == nil {
		return "no field matching key"
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

	for _, loader := range c.loaders {
		helpString := loader.HelpString(key)
		if helpString != "" {
			builder.WriteString(spaceBuffer)
			builder.WriteString(fmt.Sprintf("%s\n", helpString))
		}
	}

	return builder.String()
}

func (c *AppConfig) printHelpString() {
	fmt.Printf("%s", c.HelpString())
}
