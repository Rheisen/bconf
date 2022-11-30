package bconf

import "sync"

func NewApplicationConfiguration(appName string, appDescription string) *Config {
	return &Config{appName: appName, appDescription: appDescription}
}

type LoadCondition interface {
	Load() bool
}

type FieldCondition struct {
	FieldSetKey string
	FieldKey    string
	Condition   func(fieldValue any) bool
}

type FunctionCondition struct {
	Condition func() bool
}

type FieldSet struct {
	Key            string
	Fields         []*Field
	LoadConditions []LoadCondition
	Loaders        []Loader
}

type Config struct {
	appName        string
	appDescription string
	fieldSets      map[string]*FieldSet
	register       sync.Once
	registered     bool
}

func (c *Config) AppName() string {
	return c.appName
}

func (c *Config) AppDescription() string {
	return c.appDescription
}

func (c *Config) AddFieldSet(fieldSet *FieldSet) []error {
	return nil
}

func (c *Config) AddField(fieldSetKey string, field *Field) []error {
	return nil
}

func (c *Config) LoadFieldSet(fieldSetKey string) []error {
	return nil
}

func (c *Config) LoadField(fieldSetKey, fieldKey string) []error {
	return nil
}

func (c *Config) SetField(fieldSetKey, fieldKey string, fieldValue any) []error {
	return nil
}

// Register loads all defined field sets and optionally checks for and handles the help flag -h and --help.
func (c *Config) Register(handleHelpFlag bool) []error {
	return nil
}
