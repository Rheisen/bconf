package bconf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rheisen/bconf"
)

func TestFlagLoaderFunctions(t *testing.T) {
	l := bconf.NewFlagLoader()
	if l == nil {
		t.Fatalf("unexpected nil flag loader")
	}

	l = bconf.NewFlagLoaderWithKeyPrefix("key_prefix")
	if l.KeyPrefix != "key_prefix" {
		t.Fatalf("unexpected key prefix: %s", l.KeyPrefix)
	}
}

func TestFlagLoaderGetMap(t *testing.T) {
	const appFieldSetKey = "app"

	const idFieldKey = "id"

	const secretFieldKey = "secret"

	const appIDValue = "test-app-id"

	const appSecretValue = "test-sensitive-value"

	appIDFlagVariable := fmt.Sprintf("--%s_%s=%s", appFieldSetKey, idFieldKey, appIDValue)

	loader := bconf.NewFlagLoader()
	loader.OverrideLookup = []string{
		appIDFlagVariable,
		fmt.Sprintf("--%s_%s", appFieldSetKey, secretFieldKey),
		appSecretValue,
	}

	appMap := loader.GetMap(appFieldSetKey, []string{idFieldKey, secretFieldKey, "invalid_field_key"})
	if len(appMap) != 2 {
		t.Fatalf("unexpected appMap length '%d', expected '2'", len(appMap))
	}
}

func TestFlagLoader(t *testing.T) {
	const sessionFieldSet = "session"

	const sessionTokenKey = "token"

	const logFieldSet = "log"

	const logLevelKey = "level"

	const logColorKey = "color"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	sessionTokenFlag := fmt.Sprintf("--%s_%s", sessionFieldSet, sessionTokenKey)

	l := bconf.FlagLoader{
		OverrideLookup: []string{
			fmt.Sprintf("--%s_%s=%s", sessionFieldSet, sessionTokenKey, sessionKeyValue),
			fmt.Sprintf("--%s_%s", logFieldSet, logLevelKey),
			logLevelValue,
			fmt.Sprintf("-%s_%s", logFieldSet, logColorKey),
		},
	}
	clone := l.Clone()

	if l.Name() != "bconf_flags" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionFieldSet, sessionTokenKey), sessionTokenFlag) {
		t.Errorf("unexpected loader help string contents: '%s'", l.HelpString(sessionFieldSet, sessionTokenKey))
	}

	if !strings.Contains(clone.HelpString(sessionFieldSet, sessionTokenKey), sessionTokenFlag) {
		t.Errorf(
			"unexpected loader clone help string contents: '%s'",
			clone.HelpString(sessionFieldSet, sessionTokenKey),
		)
	}

	sessionKeyLookup, found := l.Get(sessionFieldSet, sessionTokenKey)
	if !found {
		t.Errorf("unexpected problem getting session_key value")
	}

	if sessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session key: '%s'", sessionKeyLookup)
	}

	logLevel, found := l.Get(logFieldSet, logLevelKey)
	if !found {
		t.Errorf("unexpected problem getting log_level value")
	}

	if logLevel != logLevelValue {
		t.Errorf("unexpected value for log level: '%s'", logLevel)
	}

	cloneSessionKeyLookup, cloneFound := clone.Get(sessionFieldSet, sessionTokenKey)
	if !cloneFound {
		t.Errorf("unexpected problem getting session_key value from loader clone")
	}

	if cloneSessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session_key from loader clone: '%s'", cloneSessionKeyLookup)
	}

	logColor, found := l.Get(logFieldSet, logColorKey)
	if !found {
		t.Errorf("unexpected problem getting log_color value")
	}

	if logColor != "true" {
		t.Errorf("unexpected log color value: '%s'", logColor)
	}

	_, found = l.Get("random_key", "random_key")
	if found {
		t.Errorf("not expecting to find a value for unset key")
	}
}

func TestFlagLoaderWithKeyPrefix(t *testing.T) {
	const sessionFieldSet = "session"

	const sessionTokenKey = "token"

	const logFieldSet = "log"

	const logLevelKey = "level"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	const keyPrefix = "ext_http_api"

	sessionTokenFlag := fmt.Sprintf("--%s_%s_%s", keyPrefix, sessionFieldSet, sessionTokenKey)

	l := bconf.FlagLoader{
		KeyPrefix: keyPrefix,
		OverrideLookup: []string{
			fmt.Sprintf("--%s_%s=%s", sessionFieldSet, sessionTokenKey, sessionKeyValue),
			fmt.Sprintf("--%s_%s", logFieldSet, logLevelKey),
			logLevelValue,
		},
	}
	clone := l.Clone()

	if l.Name() != "bconf_flags" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionFieldSet, sessionTokenKey), sessionTokenFlag) {
		t.Errorf("unexpected loader help string contents: '%s'", l.HelpString(sessionFieldSet, sessionTokenKey))
	}

	if !strings.Contains(clone.HelpString(sessionFieldSet, sessionTokenKey), sessionTokenFlag) {
		t.Errorf(
			"unexpected loader clone help string contents: '%s'",
			clone.HelpString(sessionFieldSet, sessionTokenKey),
		)
	}

	sessionKeyLookup, found := l.Get(sessionFieldSet, sessionTokenKey)
	if !found {
		t.Errorf("unexpected problem getting session_key value")
	}

	if sessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session key: '%s'", sessionKeyLookup)
	}

	logLevel, found := l.Get(logFieldSet, logLevelKey)
	if !found {
		t.Errorf("unexpected problem getting log_level value")
	}

	if logLevel != logLevelValue {
		t.Errorf("unexpected value for log level: '%s'", logLevel)
	}

	cloneSessionKeyLookup, cloneFound := clone.Get(sessionFieldSet, sessionTokenKey)
	if !cloneFound {
		t.Errorf("unexpected problem getting session_key value from loader clone")
	}

	if cloneSessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session_key from loader clone: '%s'", cloneSessionKeyLookup)
	}
}
