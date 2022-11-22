package bconf_test

import (
	"testing"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

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
