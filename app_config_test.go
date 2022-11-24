package bconf_test

import (
	"os"
	"testing"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

func TestNewAppConfig(t *testing.T) {
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
		Loaders:      []bconf.Loader{&bconf.EnvironmentLoader{KeyPrefix: "bconf"}},
	}

	// Phony environment

	os.Setenv("BCONF_APP_ID", "-")

	appConfig, errs := bconf.NewAppConfig(&def)
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

	t.Log(lookupValue)

	// if appID != "generate-uuid-here" {
	// 	t.Fatalf("unexpected value for appID: %s", appID)
	// }
}
