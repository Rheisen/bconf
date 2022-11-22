package bconf_test

import (
	"testing"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

func TestHelpText(t *testing.T) {
	def := baseAppConfigDefinition()

	helpText := def.HelpText()
	t.Logf("%s", helpText)
}

func TestAppConfigDefinition(t *testing.T) {
	fields := map[string]*bconf.Field{
		"app_id": {
			FieldType: bconfconst.String,
			DefaultGenerator: func() (any, error) {
				return "generated-default", nil
			},
		},
	}
	def := bconf.AppConfigDefinition{
		Name:         "app_config_test",
		ConfigFields: fields,
		KeyPrefix:    "bconf_",
		Loaders:      []string{bconfconst.EnvironmentLoader},
	}

	if errs := def.GenerateFieldDefaults(); len(errs) > 0 {
		t.Fatalf("unexpected errors generating field defaults: %v", errs)
	}

	if errs := def.Validate(); len(errs) > 0 {
		t.Fatalf("unexpected errors validating fields: %v", errs)
	}
	// appConfig, errs := bconf.NewAppConfig(def)
	// if len(errs) > 0 {
	// 	t.Fatalf("unexpected errors: %s", errs[0])
	// }

	// appID, err := appConfig.GetString("app_id")
	// if err != nil {
	// 	t.Fatalf("unexpected error getting appID: %s", err)
	// }

	// if appID != "generate-uuid-here" {
	// 	t.Fatalf("unexpected value for appID: %s", appID)
	// }
}

func baseAppConfigDefinition() bconf.AppConfigDefinition {
	fields := map[string]*bconf.Field{
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
	}
	def := bconf.AppConfigDefinition{
		Name:         "bconf_example_app",
		Description:  "Example-App is an HTTP Application for accessing weather data",
		ConfigFields: fields,
		KeyPrefix:    "bconf",
		Loaders:      []string{bconfconst.EnvironmentLoader},
	}

	return def
}
