package bconfv1_test

import (
	"os"
	"testing"

	"github.com/rheisen/bconf/bconfconst"
	"github.com/rheisen/bconf/bconfv1"
)

func TestNewAppConfig(t *testing.T) {
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

	// Phony environment

	os.Setenv("BCONF_APP_ID", "-")

	appConfig, errs := bconfv1.NewAppConfig(&def)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %s", errs[0])
	}

	appIDField, err := appConfig.GetField("app_id")
	if err != nil {
		t.Fatalf("unexpected error getting appID: %s", err)
	}

	appIDFieldValue, err := appIDField.GetValue()
	if err != nil {
		t.Fatalf("unexpected error getting appIDField value: %s", err)
	}

	t.Log(appIDFieldValue)

	lookupValue, err := appConfig.GetString("app_id")
	if err != nil {
		t.Fatalf("unexpected error getting appConfig string value: %s", err)
	}

	if lookupValue != "-" {
		t.Fatalf("unexpected app_id value, expected '%s', found '%s'", "-", lookupValue)
	}

	t.Log(lookupValue)

	// if appID != "generate-uuid-here" {
	// 	t.Fatalf("unexpected value for appID: %s", appID)
	// }
}
