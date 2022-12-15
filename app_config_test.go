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
		t.Fatalf("unexpected errors registering application configuration: %v", errs)
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

	if err = appConfig.SetField("app", "id", "user-override-value"); err != nil {
		t.Fatalf("unexpected error setting app id value: %s", err)
	}

	appID, err = appConfig.GetString("app", "id")
	if err != nil {
		t.Fatalf("unexpected error getting app_id field: %s", err)
	}

	if appID != "user-override-value" {
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

func TestAppConfigWithLoadConditions(t *testing.T) {
	const appName = "bconf_test_app"

	const appDescription = "Test-App is an HTTP server providing access to weather data"

	appConfig := bconf.NewAppConfig(
		appName,
		appDescription,
	)

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
		Fields: []*bconf.Field{loadAppOneField},
	}

	fieldSetWithLoadCondition := &bconf.FieldSet{
		Key:    "app_one",
		Fields: []*bconf.Field{},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: defaultFieldSetKey,
				FieldKey:    defaultFieldSetLoadAppOneKey,
				Condition: func(fieldValue any) bool {
					return true
				},
			},
		},
	}

	fieldSetWithUnmetLoadCondition := &bconf.FieldSet{
		Key:    "app_two",
		Fields: []*bconf.Field{},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: defaultFieldSetKey,
				FieldKey:    defaultFieldSetLoadAppTwoKey,
				Condition: func(fieldValue any) bool {
					return true
				},
			},
		},
	}

	fieldSetWithInvalidLoadCondition := &bconf.FieldSet{
		Key:    "app_three",
		Fields: []*bconf.Field{},
		LoadConditions: []bconf.LoadCondition{
			&bconf.FieldCondition{
				FieldSetKey: defaultFieldSetKey,
				Condition: func(fieldValue any) bool {
					return true
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

func TestAppConfigObservability(t *testing.T) {
	const appName = "bconf_test_app"

	const appDescription = "Test-App is an HTTP server providing access to weather data"

	appConfig := bconf.NewAppConfig(
		appName,
		appDescription,
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

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
		Fields: []*bconf.Field{idField, sessionSecretField},
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

func TestAppConfigFieldTypeGetters(t *testing.T) {
	const appName = "bconf_test_app"

	const appDescription = "Test-App is an HTTP server providing access to weather data"

	appConfig := bconf.NewAppConfig(
		appName,
		appDescription,
	)

	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

	stringsFieldKey := "strings"
	stringsFieldValue := []string{"string_one", "string_two"}
	stringsField := &bconf.Field{
		Key:       stringsFieldKey,
		FieldType: bconfconst.Strings,
		Default:   stringsFieldValue,
	}

	intsFieldKey := "ints"
	intsFieldValue := []int{1, 2}
	intsField := &bconf.Field{
		Key:       intsFieldKey,
		FieldType: bconfconst.Ints,
		Default:   intsFieldValue,
	}

	boolsFieldKey := "bools"
	boolsFieldValue := []bool{true, false}
	boolsField := &bconf.Field{
		Key:       boolsFieldKey,
		FieldType: bconfconst.Bools,
		Default:   boolsFieldValue,
	}

	timesFieldKey := "times"
	timesFieldValue := []time.Time{time.Now(), time.Now().Add(-1 * time.Hour)}
	timesField := &bconf.Field{
		Key:       timesFieldKey,
		FieldType: bconfconst.Times,
		Default:   timesFieldValue,
	}

	durationsFieldKey := "durations"
	durationsFieldValue := []time.Duration{1 * time.Minute, 1 * time.Hour}
	durationsField := &bconf.Field{
		Key:       durationsFieldKey,
		FieldType: bconfconst.Durations,
		Default:   durationsFieldValue,
	}

	defaultFieldSetKey := "default"
	defaultFieldSet := &bconf.FieldSet{
		Key: defaultFieldSetKey,
		Fields: []*bconf.Field{
			stringsField,
			intsField,
			boolsField,
			timesField,
			durationsField,
		},
	}

	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
	}

	if _, err := appConfig.GetStrings(defaultFieldSetKey, intsFieldKey); err == nil {
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

	if _, err = appConfig.GetInts(defaultFieldSetKey, stringsFieldKey); err == nil {
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

	if _, err = appConfig.GetBools(defaultFieldSetKey, stringsFieldKey); err == nil {
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

	if _, err = appConfig.GetTimes(defaultFieldSetKey, stringsFieldKey); err == nil {
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

	if _, err = appConfig.GetDurations(defaultFieldSetKey, stringsFieldKey); err == nil {
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
}
