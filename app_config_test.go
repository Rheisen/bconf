package bconf_test

import (
	"os"
	"testing"
	"time"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

func TestAppConfig(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"bconf_test_app",
		"Test-App is an HTTP server providing access to weather data",
	)

	configLoaders := []bconf.Loader{
		&bconf.EnvironmentLoader{KeyPrefix: "bconf_test"},
	}

	if errs := appConfig.SetLoaders(configLoaders...); len(errs) > 0 {
		t.Fatalf("unexpected errors setting loaders: %v", errs)
	}

	appFieldSet := &bconf.FieldSet{
		Key: "app",
		Fields: []*bconf.Field{
			{
				Key:         "id",
				FieldType:   bconfconst.String,
				Description: "Application identifier for use in application log messages and tracing",
				DefaultGenerator: func() (any, error) {
					return "generated-default", nil
				},
			},
			{
				Key:         "read_timeout",
				FieldType:   bconfconst.Duration,
				Description: "Application read timeout for HTTP requests",
				Default:     5 * time.Second,
			},
			{
				Key:       "connect_sqlite",
				FieldType: bconfconst.Bool,
				Default:   true,
			},
		},
	}

	conditionalFieldSet := &bconf.FieldSet{
		Key: "sqlite",
		Fields: []*bconf.Field{
			{
				Key:       "server",
				FieldType: bconfconst.String,
				Required:  true,
			},
		},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: "app",
				FieldKey:    "connect_sqlite",
				Condition: func(fieldValue any) bool {
					val, ok := fieldValue.(bool)
					if !ok {
						t.Fatalf("unexpected field-type value")
					}

					return val
				},
			},
		},
	}

	if errs := appConfig.AddFieldSet(conditionalFieldSet); len(errs) < 1 {
		t.Fatalf("errors expected when adding a field set with an unmet load condition dependency")
	}

	if errs := appConfig.AddFieldSet(appFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected errors adding field set: %v", errs)
	}

	if errs := appConfig.AddFieldSet(appFieldSet); len(errs) < 1 {
		t.Fatalf("errors expected when adding field set with duplicate key: %s", appFieldSet.Key)
	}

	if errs := appConfig.AddFieldSet(conditionalFieldSet); len(errs) > 1 {
		t.Fatalf("unexpected errors adding conditional field-set: %v", errs)
	}

	t.Log(appConfig.HelpString())

	if errs := appConfig.Register(false); len(errs) < 1 {
		t.Fatalf("errors expected for unset required fields")
	}

	os.Setenv("BCONF_TEST_SQLITE_SERVER", "localhost")

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected error registering application configuration: %v", errs)
	}

	appID, err := appConfig.GetString("app", "id")
	if err != nil {
		t.Fatalf("unexpected error getting app_id field: %s", err)
	}
	if appID != "generated-default" {
		t.Fatalf("unexected app_id value, found: '%s'", appID)
	}

	readTimeout, err := appConfig.GetDuration("app", "read_timeout")
	if err != nil {
		t.Fatalf("unexpected error getting app_read_timeout field: %s", err)
	}
	if readTimeout != 5*time.Second {
		t.Fatalf("unexpected app_read_timeout value, found: '%d ms'", readTimeout.Milliseconds())
	}

	os.Setenv("BCONF_TEST_APP_ID", "environment-loaded-app-id")

	appConfig.LoadField("app", "id")

	appID, err = appConfig.GetString("app", "id")
	if err != nil {
		t.Fatalf("unexpected error getting app_id field: %s", err)
	}
	if appID != "environment-loaded-app-id" {
		t.Fatalf("unexected app_id value, found: '%s'", appID)
	}
}
