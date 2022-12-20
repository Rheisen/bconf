package bconf_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

func TestAppConfigHelpString(t *testing.T) {
	appConfig := createBaseAppConfig()

	const stringFieldKey = "string"

	stringField := &bconf.Field{
		Key:         stringFieldKey,
		FieldType:   bconfconst.String,
		Description: "string field description",
		DefaultGenerator: func() (any, error) {
			return "some_value", nil
		},
		Enumeration: []any{"some_value", "other_value", "another_value"},
	}

	const secretStringFieldKey = "string_secret"

	secretStringField := &bconf.Field{
		Key:       secretStringFieldKey,
		FieldType: bconfconst.String,
		Sensitive: true,
		Default:   "some-super-secret-value",
		Validator: func(fieldValue any) error {
			return nil
		},
	}

	const intFieldKey = "int"

	intField := &bconf.Field{
		Key:       intFieldKey,
		FieldType: bconfconst.Int,
		Required:  true,
	}

	const durationFieldKey = "duration"

	durationField := &bconf.Field{
		Key:       durationFieldKey,
		FieldType: bconfconst.Duration,
		Required:  true,
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key:    defaultFieldSetKey,
		Fields: bconf.Fields{stringField, secretStringField, intField},
	}

	const conditionalFieldSetKey = "conditional"

	conditionalFieldSet := &bconf.FieldSet{
		Key:    conditionalFieldSetKey,
		Fields: bconf.Fields{durationField},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: defaultFieldSetKey,
				FieldKey:    stringFieldKey,
				Condition: func(fieldValue any) (bool, error) {
					val, ok := fieldValue.(string)
					if !ok {
						return false, fmt.Errorf("unexpected field value type")
					}

					return val == "some_value", nil
				},
			},
		},
	}

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected errors adding field-set: %v", errs)
	}

	if errs := appConfig.AddFieldSet(conditionalFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected errors adding conditional field-set: %v", errs)
	}

	t.Log(appConfig.HelpString())
}

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
		Fields: bconf.Fields{
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
		Fields: bconf.Fields{
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
				Condition: func(fieldValue any) (bool, error) {
					val, ok := fieldValue.(bool)
					if !ok {
						return false, fmt.Errorf("unexpected field-type value")
					}

					return val, nil
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

	if errs := appConfig.Register(false); len(errs) < 1 {
		t.Fatalf("errors expected for unset required fields")
	}

	os.Setenv("BCONF_TEST_SQLITE_SERVER", "localhost")

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected errors registering application configuration: %v", errs)
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
		Fields: bconf.Fields{
			idFieldInvalidDefaultGenerator,
			readTimeoutFieldInvalidDefault,
		},
	}

	if errs := appConfig.AddFieldSet(invalidAppFieldSet); len(errs) < 1 {
		t.Fatalf("expected errors adding field set with invalid fields")
	}

	fieldSetWithEmptyField := &bconf.FieldSet{
		Key: "default",
		Fields: bconf.Fields{
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
		Fields: bconf.Fields{fieldWithDefaultAndRequiredSet},
	}

	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
		t.Fatalf("expected an error adding field with default and required set")
	}

	fieldSetWithInvalidField.Fields = bconf.Fields{fieldWithDefaultNotInEnumeration}

	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
		t.Fatalf("expected an error adding field with default value not in enumeration")
	}

	fieldSetWithInvalidField.Fields = bconf.Fields{fieldWithGeneratedDefaultNotInEnumeration}

	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
		t.Fatalf("expected an error adding field with generated default value not in enumeration")
	}
}

func TestAppConfigWithLoadConditions(t *testing.T) {
	appConfig := createBaseAppConfig()

	const defaultFieldSetKey = "default"

	const defaultFieldSetLoadAppOneKey = "load_app_one"

	const defaultFieldSetLoadAppTwoKey = "load_app_two"

	loadAppOneField := &bconf.Field{
		Key:       defaultFieldSetLoadAppOneKey,
		FieldType: bconfconst.Bool,
		Required:  true,
	}

	defaultFieldSet := &bconf.FieldSet{
		Key:    "default",
		Fields: bconf.Fields{loadAppOneField},
	}

	fieldSetWithLoadCondition := &bconf.FieldSet{
		Key:    "app_one",
		Fields: bconf.Fields{},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: defaultFieldSetKey,
				FieldKey:    defaultFieldSetLoadAppOneKey,
				Condition: func(fieldValue any) (bool, error) {
					return true, nil
				},
			},
		},
	}

	fieldSetWithUnmetLoadCondition := &bconf.FieldSet{
		Key:    "app_two",
		Fields: bconf.Fields{},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: defaultFieldSetKey,
				FieldKey:    defaultFieldSetLoadAppTwoKey,
				Condition: func(fieldValue any) (bool, error) {
					return true, nil
				},
			},
		},
	}

	fieldSetWithInvalidLoadCondition := &bconf.FieldSet{
		Key:    "app_three",
		Fields: bconf.Fields{},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: defaultFieldSetKey,
				Condition: func(fieldValue any) (bool, error) {
					return true, nil
				},
			},
		},
	}

	if errs := appConfig.AddFieldSet(fieldSetWithLoadCondition); len(errs) < 1 {
		t.Fatalf("expected error adding field set with unmet field-set load condition")
	}

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
	}

	if errs := appConfig.AddFieldSet(fieldSetWithUnmetLoadCondition); len(errs) < 1 {
		t.Fatalf("expected error adding field set with unmet field load condition")
	}

	if errs := appConfig.AddFieldSet(fieldSetWithLoadCondition); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding field set with valid load condition: %v", errs)
	}

	if errs := appConfig.AddFieldSet(fieldSetWithInvalidLoadCondition); len(errs) < 1 {
		t.Fatalf("expected error adding field set with invalid load condition")
	}
}

func TestAppConfigAddFieldSets(t *testing.T) {
	appConfig := createBaseAppConfig()

	fieldSetOne := &bconf.FieldSet{
		Key:    "one",
		Fields: bconf.Fields{},
	}
	fieldSetTwo := &bconf.FieldSet{
		Key:    "two",
		Fields: bconf.Fields{},
	}
	fieldSetThree := &bconf.FieldSet{
		Key:    "three",
		Fields: bconf.Fields{},
	}
	fieldSetFour := &bconf.FieldSet{
		Fields: bconf.Fields{},
	}

	if errs := appConfig.AddFieldSets(fieldSetOne, fieldSetTwo); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding field-sets: %v", errs)
	}

	if errs := appConfig.AddFieldSets(fieldSetThree, fieldSetFour); len(errs) < 1 {
		t.Fatalf("expected an error adding field-set with missing key")
	} else if !strings.Contains(errs[0].Error(), "field-set key required") {
		t.Fatalf("unexpected error message: %s", errs[0])
	}

	if keys := appConfig.GetFieldSetKeys(); len(keys) != 2 {
		t.Fatalf("unexpected number of field-sets found on app config: %d", len(keys))
	}
}

func TestAppConfigAddField(t *testing.T) {
	appConfig := createBaseAppConfig()

	fieldSetOne := &bconf.FieldSet{
		Key:    "one",
		Fields: bconf.Fields{},
	}

	idFieldKey := "id"
	idFieldGeneratedDefaultValue := "generated-default-value"
	idField := &bconf.Field{
		Key:         idFieldKey,
		FieldType:   bconfconst.String,
		Description: "Application identifier for use in application log messages and tracing",
		DefaultGenerator: func() (any, error) {
			return idFieldGeneratedDefaultValue, nil
		},
	}

	fieldWithGenerateDefaultError := &bconf.Field{
		Key:       "field_generate_default_error",
		FieldType: bconfconst.String,
		DefaultGenerator: func() (any, error) {
			return "", errors.New("generated error")
		},
	}

	fieldMissingFieldType := &bconf.Field{
		Key: "field_missing_field_type",
	}

	if errs := appConfig.AddFieldSets(fieldSetOne); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding field-sets: %v", errs)
	}

	if errs := appConfig.AddField("one", idField); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding field: %v", errs)
	}

	if errs := appConfig.AddField("one", idField); len(errs) < 1 {
		t.Fatalf("expected error trying to add duplicate field to field-set")
	}

	if errs := appConfig.AddField("undefined_field_set_key", idField); len(errs) < 1 {
		t.Fatalf("expected error trying to add field to undefined field-set")
	}

	if errs := appConfig.AddField("one", fieldWithGenerateDefaultError); len(errs) < 1 {
		t.Fatalf("expected error trying to add field with bad generated default")
	}

	if errs := appConfig.AddField("one", fieldMissingFieldType); len(errs) < 1 {
		t.Fatalf("expected error trying to add field with missing field-type")
	}
}

func TestAppConfigObservability(t *testing.T) {
	appConfig := createBaseAppConfig()

	idFieldKey := "id"
	idFieldGeneratedDefaultValue := "generated-default-value"
	idField := &bconf.Field{
		Key:         idFieldKey,
		FieldType:   bconfconst.String,
		Description: "Application identifier for use in application log messages and tracing",
		DefaultGenerator: func() (any, error) {
			return idFieldGeneratedDefaultValue, nil
		},
	}

	sessionSecretFieldKey := "session_secret"
	sessionSecretEnvironmentValue := "environment-session-secret-value"
	sessionSecretField := &bconf.Field{
		Key:       sessionSecretFieldKey,
		FieldType: bconfconst.String,
		Sensitive: true,
		Validator: func(fieldValue any) error {
			return nil
		},
	}

	os.Setenv("APP_SESSION_SECRET", sessionSecretEnvironmentValue)

	appFieldSetKey := "app"
	appFieldSet := &bconf.FieldSet{
		Key:    appFieldSetKey,
		Fields: bconf.Fields{idField, sessionSecretField},
	}

	if errs := appConfig.AddFieldSet(appFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected errors adding app field set: %v", errs)
	}

	foundFieldSetKeys := appConfig.GetFieldSetKeys()
	if len(foundFieldSetKeys) != 1 {
		t.Fatalf("unexpected length of field set keys returned from app config: %d", len(foundFieldSetKeys))
	}

	if foundFieldSetKeys[0] != appFieldSetKey {
		t.Fatalf("unexpected field-set key in keys returned from app config: '%s'", foundFieldSetKeys[0])
	}

	foundAppFieldSetKeys, err := appConfig.GetFieldSetFieldKeys(appFieldSetKey)
	if err != nil {
		t.Fatalf("unexpected issue getting app field-set field keys: %s", err)
	}

	if len(foundAppFieldSetKeys) < len(appFieldSet.Fields) {
		t.Fatalf("length of field-set field keys does not match the length of fields: %d", len(foundAppFieldSetKeys))
	}

	fieldMap := appConfig.ConfigMap()

	if _, found := fieldMap[appFieldSetKey]; !found {
		t.Fatalf("expected to find app field-set key in config map")
	}

	if _, found := fieldMap[appFieldSetKey][idFieldKey]; !found {
		t.Fatalf("expected to find app id key in config map")
	}

	if _, found := fieldMap[appFieldSetKey][sessionSecretFieldKey]; found {
		t.Fatalf("unexpected session-secret key found in config map when no value is set")
	}

	if errs := appConfig.Register(false); errs != nil {
		t.Fatalf("unexpected errors registering app config: %v", errs)
	}

	fieldMap = appConfig.ConfigMap()

	fieldMapValue, found := fieldMap[appFieldSetKey][sessionSecretFieldKey]
	if !found {
		t.Fatalf("expected to find session-secret key in config map")
	}

	if fieldMapValue == sessionSecretEnvironmentValue {
		t.Fatalf(
			"unexpected sensitive value (%s) output in config map values: '%s'",
			sessionSecretFieldKey,
			fieldMapValue,
		)
	}
}

func TestAppConfigSetField(t *testing.T) {
	appConfig := createBaseAppConfig()

	const stringFieldKey = "string"

	const stringFieldValue = "string_one"

	stringField := &bconf.Field{
		Key:         stringFieldKey,
		FieldType:   bconfconst.String,
		Default:     stringFieldValue,
		Enumeration: []any{"string_one", "string_two", "string_three"},
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			stringField,
		},
	}

	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, "some_val"); err == nil {
		t.Fatalf("expected error setting field when field-set is not present")
	} else if !strings.Contains(err.Error(), fmt.Sprintf("field-set with key '%s' not found", defaultFieldSetKey)) {
		t.Fatalf("unexpected error message: %s", err.Error())
	}

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
	}

	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, 3928482); err == nil {
		t.Fatalf("expected error setting field to mismatched type")
	} else if !strings.Contains(err.Error(), "invalid value field-type") {
		t.Fatalf("unexpected error message when setting field to mismatched field-type: %s", err)
	}

	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, "string_zero"); err == nil {
		t.Fatalf("expected error setting field to value not in enumeration list")
	} else if !strings.Contains(err.Error(), "value not found in enumeration list") {
		t.Fatalf("unexpected error message when setting field to value not in enumeraiton list: %s", err)
	}

	if err := appConfig.SetField(defaultFieldSetKey, "some_key", "some_val"); err == nil {
		t.Fatalf("expected error setting field when field is not present")
	} else if !strings.Contains(err.Error(), "field with key") {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestAppConfigReloadingFields(t *testing.T) {
	appConfig := createBaseAppConfig()

	const stringFieldKey = "string"

	const stringFieldValue = "string_one"

	stringField := &bconf.Field{
		Key:       stringFieldKey,
		FieldType: bconfconst.String,
		Default:   stringFieldValue,
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			stringField,
		},
	}

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
	}

	if errs := appConfig.LoadFieldSet(defaultFieldSetKey); len(errs) != 1 {
		t.Fatalf("expected error loading field-set before the app-config is registered")
	} else if !strings.Contains(errs[0].Error(), "cannot be called before the app-config has been registered") {
		t.Fatalf("unexpected error message when loading field-set before app-config is registered: %s", errs[0])
	}

	if errs := appConfig.LoadField(defaultFieldSetKey, stringFieldKey); len(errs) != 1 {
		t.Fatalf("expected error loading field before the app-config is registered")
	} else if !strings.Contains(errs[0].Error(), "cannot be called before the app-config has been registered") {
		t.Fatalf("unexpected error message when loading field before app-config is registered: %s", errs[0])
	}

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected error(s) registering app-config: %v", errs)
	}

	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, stringFieldKey)), "string_two")

	if errs := appConfig.LoadFieldSet(defaultFieldSetKey); len(errs) > 0 {
		t.Fatalf("unexpected errors loading field-set: %v", errs)
	}

	if val, err := appConfig.GetString(defaultFieldSetKey, stringFieldKey); err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if val != "string_two" {
		t.Fatalf("unexpected field value: '%s'", val)
	}

	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, stringFieldKey)), "string_three")

	if errs := appConfig.LoadField(defaultFieldSetKey, stringFieldKey); len(errs) > 0 {
		t.Fatalf("unexpected errors loading field: %v", errs)
	}

	if val, err := appConfig.GetString(defaultFieldSetKey, stringFieldKey); err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if val != "string_three" {
		t.Fatalf("unexpected field value: '%s'", val)
	}
}

func TestAppConfigFieldValidators(t *testing.T) {
	appConfig := createBaseAppConfig()

	const stringFieldKey = "string"

	stringFieldValue := "string_one"
	validatorExpectedValue := "string_two"
	validatorErrorString := fmt.Sprintf("expected value to be '%s'", validatorExpectedValue)
	stringField := &bconf.Field{
		Key:       stringFieldKey,
		FieldType: bconfconst.String,
		Default:   stringFieldValue,
		Validator: func(fieldValue any) error {
			val, _ := fieldValue.(string)

			if val != validatorExpectedValue {
				return errors.New(validatorErrorString)
			}

			return nil
		},
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			stringField,
		},
	}

	expectContains := fmt.Sprintf(
		"invalid default value: error from field validator: %s",
		validatorErrorString,
	)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) != 1 {
		t.Fatalf("expected 1 error adding default field-set with default value not passing validator: %v", errs)
	} else if !strings.Contains(errs[0].Error(), expectContains) {
		t.Fatalf("unexpected error message: %s", errs[0])
	}

	stringField.Default = nil
	stringField.DefaultGenerator = func() (any, error) {
		return stringFieldValue, nil
	}

	expectContains = fmt.Sprintf(
		"invalid generated default value: error from field validator: %s",
		validatorErrorString,
	)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) != 1 {
		t.Fatalf(
			"expected 1 error adding default field-set with generated default value not passing validator: %v",
			errs,
		)
	} else if !strings.Contains(errs[0].Error(), expectContains) {
		t.Fatalf("unexpected error message: %s", errs[0])
	}

	stringField.Default = validatorExpectedValue
	stringField.DefaultGenerator = nil

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
	}

	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, stringFieldValue); err == nil {
		t.Fatalf("expected error setting field value violating validator func")
	}
}

func TestAppConfigFieldDefaultGenerators(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"app",
		"description",
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	const stringFieldKey = "string"

	defaultGeneratorError := "problem generating default"
	stringField := &bconf.Field{
		Key:       stringFieldKey,
		FieldType: bconfconst.String,
		DefaultGenerator: func() (any, error) {
			return nil, errors.New(defaultGeneratorError)
		},
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			stringField,
		},
	}

	expectContains := fmt.Sprintf(
		"default value generation error: problem generating default field value: %s",
		defaultGeneratorError,
	)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) != 1 {
		t.Fatalf(
			"expected 1 error adding default field-set with generated default value function error: %v",
			errs,
		)
	} else if !strings.Contains(errs[0].Error(), expectContains) {
		t.Fatalf("unexpected error message: %s", errs[0])
	}
}

func TestAppConfigStringFieldTypes(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"app",
		"description",
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	stringsFieldKey := "strings"
	stringsFieldValue := []string{"string_one", "string_two"}
	stringsEnvValue := "string_three, string_four"
	stringsParsedEnvValue := []string{"string_three", "string_four"}
	stringsField := &bconf.Field{
		Key:       stringsFieldKey,
		FieldType: bconfconst.Strings,
		Default:   stringsFieldValue,
	}

	herringFieldKey := "ints"
	herringField := &bconf.Field{
		Key:       herringFieldKey,
		FieldType: bconfconst.Ints,
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			stringsField,
			herringField,
		},
	}

	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, stringsFieldKey)), stringsEnvValue)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
	}

	if _, err := appConfig.GetStrings(defaultFieldSetKey, herringFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundStringVals, err := appConfig.GetStrings(defaultFieldSetKey, stringsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundStringVals {
		if stringsFieldValue[idx] != val {
			t.Errorf("unexpected value found: '%s', expected '%s", val, stringsFieldValue[idx])
		}
	}

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected error(s) registering app config: %v", errs)
	}

	foundStringVals, err = appConfig.GetStrings(defaultFieldSetKey, stringsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundStringVals {
		if stringsParsedEnvValue[idx] != val {
			t.Errorf("unexpected value found: '%s', expected '%s", val, stringsParsedEnvValue[idx])
		}
	}
}

func TestAppConfigIntFieldTypes(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"app",
		"description",
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	intFieldKey := "int"
	intFieldValue := 1
	intEnvValue := "2"
	intParsedEnvValue := 2
	intField := &bconf.Field{
		Key:       intFieldKey,
		FieldType: bconfconst.Int,
		Default:   intFieldValue,
	}

	intsFieldKey := "ints"
	intsFieldValue := []int{1, 2}
	intsEnvValue := "3, 4"
	intsParsedEnvValue := []int{3, 4}
	intsField := &bconf.Field{
		Key:       intsFieldKey,
		FieldType: bconfconst.Ints,
		Default:   intsFieldValue,
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			intField,
			intsField,
		},
	}

	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, intFieldKey)), intEnvValue)
	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, intsFieldKey)), intsEnvValue)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
	}

	if _, err := appConfig.GetInt(defaultFieldSetKey, intsFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundIntVal, err := appConfig.GetInt(defaultFieldSetKey, intFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundIntVal != intFieldValue {
		t.Errorf("unexpected value found: '%d', expected '%d", foundIntVal, intFieldValue)
	}

	if _, err = appConfig.GetInts(defaultFieldSetKey, intFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundIntVals, err := appConfig.GetInts(defaultFieldSetKey, intsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundIntVals {
		if intsFieldValue[idx] != val {
			t.Errorf("unexpected value found: '%d', expected '%d", val, intsFieldValue[idx])
		}
	}

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected error(s) registering app config: %v", errs)
	}

	foundIntVal, err = appConfig.GetInt(defaultFieldSetKey, intFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundIntVal != intParsedEnvValue {
		t.Errorf("unexpected value found: '%d', expected '%d", foundIntVal, intParsedEnvValue)
	}

	foundIntVals, err = appConfig.GetInts(defaultFieldSetKey, intsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundIntVals {
		if intsParsedEnvValue[idx] != val {
			t.Errorf("unexpected value found: '%d', expected '%d", val, intsParsedEnvValue[idx])
		}
	}
}

func TestAppConfigBoolFieldTypes(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"app",
		"description",
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	boolFieldKey := "bool"
	boolFieldValue := true
	boolEnvValue := "false"
	boolParsedEnvValue := false
	boolField := &bconf.Field{
		Key:       boolFieldKey,
		FieldType: bconfconst.Bool,
		Default:   boolFieldValue,
	}

	boolsFieldKey := "bools"
	boolsFieldValue := []bool{true, false}
	boolsEnvValue := "false, true"
	boolsParsedEnvValue := []bool{false, true}
	boolsField := &bconf.Field{
		Key:       boolsFieldKey,
		FieldType: bconfconst.Bools,
		Default:   boolsFieldValue,
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			boolField,
			boolsField,
		},
	}

	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, boolFieldKey)), boolEnvValue)
	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, boolsFieldKey)), boolsEnvValue)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
	}

	if _, err := appConfig.GetBool(defaultFieldSetKey, boolsFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundBoolVal, err := appConfig.GetBool(defaultFieldSetKey, boolFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundBoolVal != boolFieldValue {
		t.Errorf("unexpected value found: '%v', expected '%v", foundBoolVal, boolFieldValue)
	}

	if _, err = appConfig.GetBools(defaultFieldSetKey, boolFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundBoolVals, err := appConfig.GetBools(defaultFieldSetKey, boolsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundBoolVals {
		if boolsFieldValue[idx] != val {
			t.Errorf("unexpected value found: '%v', expected '%v", val, boolsFieldValue[idx])
		}
	}

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected error(s) registering app config: %v", errs)
	}

	foundBoolVal, err = appConfig.GetBool(defaultFieldSetKey, boolFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundBoolVal != boolParsedEnvValue {
		t.Errorf("unexpected value found: '%v', expected '%v", foundBoolVal, boolParsedEnvValue)
	}

	foundBoolVals, err = appConfig.GetBools(defaultFieldSetKey, boolsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundBoolVals {
		if boolsParsedEnvValue[idx] != val {
			t.Errorf("unexpected value found: '%v', expected '%v", val, boolsParsedEnvValue[idx])
		}
	}
}

func TestAppConfigDurationFieldTypes(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"app",
		"description",
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	durationFieldKey := "duration"
	durationFieldValue := 1 * time.Minute
	durationEnvValue := "1h"
	durationParsedEnvValue := 1 * time.Hour
	durationField := &bconf.Field{
		Key:       durationFieldKey,
		FieldType: bconfconst.Duration,
		Default:   durationFieldValue,
	}

	durationsFieldKey := "durations"
	durationsFieldValue := []time.Duration{1 * time.Minute, 1 * time.Hour}
	durationsEnvValue := "1h, 1m"
	durationsParsedEnvValue := []time.Duration{1 * time.Hour, 1 * time.Minute}
	durationsField := &bconf.Field{
		Key:       durationsFieldKey,
		FieldType: bconfconst.Durations,
		Default:   durationsFieldValue,
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			durationField,
			durationsField,
		},
	}

	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, durationFieldKey)), durationEnvValue)
	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, durationsFieldKey)), durationsEnvValue)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
	}

	if _, err := appConfig.GetDuration(defaultFieldSetKey, durationsFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundDurationVal, err := appConfig.GetDuration(defaultFieldSetKey, durationFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundDurationVal != durationFieldValue {
		t.Errorf("unexpected value found: '%s', expected '%s", foundDurationVal.String(), durationFieldValue.String())
	}

	if _, err = appConfig.GetDurations(defaultFieldSetKey, durationFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundDurationVals, err := appConfig.GetDurations(defaultFieldSetKey, durationsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundDurationVals {
		if durationsFieldValue[idx] != val {
			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), durationsFieldValue[idx].String())
		}
	}

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected error(s) registering app config: %v", errs)
	}

	foundDurationVal, err = appConfig.GetDuration(defaultFieldSetKey, durationFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundDurationVal != durationParsedEnvValue {
		t.Errorf(
			"unexpected value found: '%s', expected '%s",
			foundDurationVal.String(),
			durationParsedEnvValue.String(),
		)
	}

	foundDurationVals, err = appConfig.GetDurations(defaultFieldSetKey, durationsFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundDurationVals {
		if durationsParsedEnvValue[idx] != val {
			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), durationsParsedEnvValue[idx].String())
		}
	}
}

func TestAppConfigTimeFieldTypes(t *testing.T) {
	appConfig := bconf.NewAppConfig(
		"app",
		"description",
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	baseTime := time.Now()

	timeFieldKey := "time"
	timeFieldValue := baseTime
	timeEnvValue := baseTime.Add(-1 * time.Hour).Format(time.RFC3339)
	timeParsedEnvValue := baseTime.Add(-1 * time.Hour)
	timeField := &bconf.Field{
		Key:       timeFieldKey,
		FieldType: bconfconst.Time,
		Default:   timeFieldValue,
	}

	timesFieldKey := "times"
	timesFieldValue := []time.Time{baseTime, baseTime.Add(-1 * time.Hour)}
	timesEnvValue := fmt.Sprintf(
		"%s, %s",
		baseTime.Add(-1*time.Hour).Format(time.RFC3339),
		baseTime.Format(time.RFC3339),
	)
	timesParsedEnvValue := []time.Time{baseTime.Add(-1 * time.Hour), baseTime}
	timesField := &bconf.Field{
		Key:       timesFieldKey,
		FieldType: bconfconst.Times,
		Default:   timesFieldValue,
	}

	const defaultFieldSetKey = "default"

	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: bconf.Fields{
			timeField,
			timesField,
		},
	}

	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, timeFieldKey)), timeEnvValue)
	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, timesFieldKey)), timesEnvValue)

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
	}

	if _, err := appConfig.GetTime(defaultFieldSetKey, timesFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundTimeVal, err := appConfig.GetTime(defaultFieldSetKey, timeFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundTimeVal != timeFieldValue {
		t.Errorf("unexpected value found: '%s', expected '%s", foundTimeVal.String(), timeFieldValue.String())
	}

	if _, err = appConfig.GetTimes(defaultFieldSetKey, timeFieldKey); err == nil {
		t.Fatalf("expected error getting mismatched field type")
	}

	foundTimeVals, err := appConfig.GetTimes(defaultFieldSetKey, timesFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundTimeVals {
		if !timesFieldValue[idx].Equal(val) {
			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), timesFieldValue[idx].String())
		}
	}

	if errs := appConfig.Register(false); len(errs) > 0 {
		t.Fatalf("unexpected error(s) registering app config: %v", errs)
	}

	foundTimeVal, err = appConfig.GetTime(defaultFieldSetKey, timeFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	} else if foundTimeVal.Format(time.RFC3339) != timeParsedEnvValue.Format(time.RFC3339) {
		t.Errorf("unexpected value found: '%s', expected '%s", foundTimeVal.String(), timeParsedEnvValue.String())
	}

	foundTimeVals, err = appConfig.GetTimes(defaultFieldSetKey, timesFieldKey)
	if err != nil {
		t.Fatalf("unexpected error getting field value: %s", err)
	}

	for idx, val := range foundTimeVals {
		if timesParsedEnvValue[idx].Format(time.RFC3339) != val.Format(time.RFC3339) {
			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), timesParsedEnvValue[idx].String())
		}
	}
}

func createBaseAppConfig() *bconf.AppConfig {
	appConfig := bconf.NewAppConfig(
		"app",
		"description",
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	return appConfig
}
