package bconfv1_test

import (
	"testing"

	"github.com/rheisen/bconf/bconfconst"
	"github.com/rheisen/bconf/bconfv1"
)

func TestHelpText(t *testing.T) {
	def := baseAppConfigDefinition()

	helpText := def.HelpText()
	expectedHelpText := `Usage of 'bconf_example_app':
Example-App is an HTTP Application for accessing weather data

Required Configuration:
	session_secret string
		Session secret used for user authentication
		Environment key: 'BCONF_SESSION_SECRET'
Optional Configuration:
	app_id string
		Application identifier for use in application log messages and tracing
		Default value: <generated-at-run-time>
		Environment key: 'BCONF_APP_ID'
	log_color bool
		Log in color when console logging is set
		Default value: 'true'
		Environment key: 'BCONF_LOG_COLOR'
	log_config string
		Logging configuration presets
		Accepted values: ['production', 'development']
		Default value: 'production'
		Environment key: 'BCONF_LOG_CONFIG'
`
	if helpText != expectedHelpText {
		t.Fatalf("unexpected help text value, expected:%sfound:%s", expectedHelpText, helpText)
	}
}

func TestClone(t *testing.T) {
	def := baseAppConfigDefinition()
	clone := def.Clone()

	newName := "New definition name"
	newDescription := "New definition description"
	newAppIDDescription := "New app_id description"

	// Test editing definition fields
	def.Name = newName
	def.Description = newDescription

	if clone.Name == newName {
		t.Fatalf("unexpected 'Name' value in cloned app config definition: %s", clone.Name)
	}
	if clone.Description == newDescription {
		t.Fatalf("unexpected 'Description' value in cloned app config definition: %s", clone.Description)
	}

	// Test editing a config field
	def.ConfigFields["app_id"].Description = newAppIDDescription
	if def.ConfigFields["app_id"].Description != newAppIDDescription {
		t.Fatalf("problem setting app config definition 'app_id' description")
	}
	if clone.ConfigFields["app_id"].Description == "New app_id description" {
		t.Fatalf("unexpected value in cloned app config definition: %s", clone.ConfigFields["app_id"].Description)
	}

	// Test deleting a config fields map value
	delete(def.ConfigFields, "app_id")
	if _, found := def.ConfigFields["app_id"]; found {
		t.Fatalf("problem deleting app config definition entry")
	}

	if _, found := clone.ConfigFields["app_id"]; !found {
		t.Fatalf("problem finding 'app_id' field on app config definition clone after deleting from original")
	}
}

func TestAppConfigDefinition(t *testing.T) {
	fields := map[string]*bconfv1.Field{
		"app_id": {
			FieldType: bconfconst.String,
			DefaultGenerator: func() (any, error) {
				return "generated-default", nil
			},
		},
	}
	def := bconfv1.AppConfigDefinition{
		Name:         "app_config_test",
		ConfigFields: fields,
		Loaders:      []bconfv1.Loader{&bconfv1.EnvironmentLoader{KeyPrefix: "bconf"}},
	}

	if errs := def.GenerateFieldDefaults(); len(errs) > 0 {
		t.Fatalf("unexpected errors generating field defaults: %v", errs)
	}

	if errs := def.Validate(); len(errs) > 0 {
		t.Fatalf("unexpected errors validating fields: %v", errs)
	}
}

func baseAppConfigDefinition() bconfv1.AppConfigDefinition {
	fields := map[string]*bconfv1.Field{
		"app_id": {
			FieldType:   bconfconst.String,
			Description: "Application identifier for use in application log messages and tracing",
			DefaultGenerator: func() (any, error) {
				return "generated-default", nil
			},
		},
		"session_secret": {
			FieldType:   bconfconst.String,
			Description: "Session secret used for user authentication",
			Required:    true,
		},
		"log_config": {
			FieldType:   bconfconst.String,
			Description: "Logging configuration presets",
			Default:     "production",
			Enumeration: []any{"production", "development"},
		},
		"log_color": {
			FieldType:   bconfconst.Bool,
			Description: "Log in color when console logging is set",
			Default:     true,
		},
	}
	def := bconfv1.AppConfigDefinition{
		Name:         "bconf_example_app",
		Description:  "Example-App is an HTTP Application for accessing weather data",
		ConfigFields: fields,
		Loaders:      []bconfv1.Loader{&bconfv1.EnvironmentLoader{KeyPrefix: "bconf"}},
	}

	return def
}
