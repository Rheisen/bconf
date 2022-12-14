package bconf_test

import (
	"os"
	"testing"
	"time"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

func TestAppConfig(t *testing.T) {
	const appName = "bconf_test_app"

	const appDescription = "Test-App is an HTTP server providing access to weather data"

	appConfig := bconf.NewAppConfig(
		appName,
		appDescription,
	)

	if appConfig.AppName() != appName {
		t.Errorf("unexpected value returned from AppName(): '%s'", appConfig.AppName())
	}

	if appConfig.AppDescription() != appDescription {
		t.Errorf("unexpected value returned from AppDescription(): '%s'", appConfig.AppDescription())
	}

	configLoaders := []bconf.Loader{
		&bconf.EnvironmentLoader{KeyPrefix: "bconf_test"},
	}

	if errs := appConfig.SetLoaders(configLoaders...); len(errs) > 0 {
		t.Fatalf("unexpected errors setting loaders: %v", errs)
	}

	const appGeneratedID = "generated-default"

	appFieldSet := &bconf.FieldSet{
		Key: "app",
		Fields: []*bconf.Field{
			{
				Key:         "id",
				FieldType:   bconfconst.String,
				Description: "Application identifier for use in application log messages and tracing",
				DefaultGenerator: func() (any, error) {
					return appGeneratedID, nil
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

	if appID != appGeneratedID {
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

func TestBadAppConfigFields(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"bconf_test_app",
		"Test-App is an HTTP server providing access to weather data",
	)

	configLoadersWithDuplicates := []bconf.Loader{
		&bconf.EnvironmentLoader{KeyPrefix: "bconf_test"},
		&bconf.EnvironmentLoader{},
	}

	if errs := appConfig.SetLoaders(configLoadersWithDuplicates...); len(errs) < 1 {
		t.Fatalf("expected error setting loaders with duplicates")
	}

	configLoaders := []bconf.Loader{
		&bconf.EnvironmentLoader{KeyPrefix: "bconf_test"},
	}

	if errs := appConfig.SetLoaders(configLoaders...); len(errs) > 0 {
		t.Fatalf("unexpected errors setting loaders: %v", errs)
	}

	idFieldInvalidDefaultGenerator := &bconf.Field{
		Key:         "id",
		FieldType:   bconfconst.Int,
		Description: "Application identifier for use in application log messages and tracing",
		DefaultGenerator: func() (any, error) {
			return "generated-default", nil
		},
	}
	readTimeoutFieldInvalidDefault := &bconf.Field{
		Key:         "read_timeout",
		FieldType:   bconfconst.Duration,
		Description: "Application read timeout for HTTP requests",
		Default:     5,
	}
	emptyFieldSet := &bconf.FieldSet{}

	if errs := appConfig.AddFieldSet(emptyFieldSet); len(errs) < 1 {
		t.Fatalf("expected error adding empty field set")
	}

	invalidAppFieldSet := &bconf.FieldSet{
		Key: "app",
		Fields: []*bconf.Field{
			idFieldInvalidDefaultGenerator,
			readTimeoutFieldInvalidDefault,
		},
	}

	if errs := appConfig.AddFieldSet(invalidAppFieldSet); len(errs) < 1 {
		t.Fatalf("expected errors adding field set with invalid fields")
	}

	fieldSetWithEmptyField := &bconf.FieldSet{
		Key: "default",
		Fields: []*bconf.Field{
			{},
		},
	}

	if errs := appConfig.AddFieldSet(fieldSetWithEmptyField); len(errs) < 2 {
		t.Fatalf("expected at least two errors adding a field-set with an empty field")
	}

	fieldWithDefaultAndRequiredSet := &bconf.Field{
		Key:       "log_level",
		FieldType: bconfconst.String,
		Default:   "info",
		Required:  true,
	}

	fieldWithDefaultNotInEnumeration := &bconf.Field{
		Key:         "log_level",
		FieldType:   bconfconst.String,
		Default:     "fatal",
		Enumeration: []any{"debug", "info", "warn", "error"},
	}

	fieldWithGeneratedDefaultNotInEnumeration := &bconf.Field{
		Key:       "log_level",
		FieldType: bconfconst.String,
		DefaultGenerator: func() (any, error) {
			return "fatal", nil
		},
		Enumeration: []any{"debug", "info", "warn", "error"},
	}

	fieldSetWithInvalidField := &bconf.FieldSet{
		Key:    "default",
		Fields: []*bconf.Field{fieldWithDefaultAndRequiredSet},
	}

	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
		t.Fatalf("expected an error adding field with default and required set")
	}

	fieldSetWithInvalidField.Fields = []*bconf.Field{fieldWithDefaultNotInEnumeration}

	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
		t.Fatalf("expected an error adding field with default value not in enumeration")
	}

	fieldSetWithInvalidField.Fields = []*bconf.Field{fieldWithGeneratedDefaultNotInEnumeration}

	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
		t.Fatalf("expected an error adding field with generated default value not in enumeration")
	}
}
