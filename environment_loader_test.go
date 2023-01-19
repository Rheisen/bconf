package bconf_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rheisen/bconf"
)

func TestEnvironmentLoaderFunctions(t *testing.T) {
	l := bconf.NewEnvironmentLoader()
	if l == nil {
		t.Fatalf("unexpected nil environment loader")
	}

	l = bconf.NewEnvironmentLoaderWithKeyPrefix("key_prefix")
	if l.KeyPrefix != "key_prefix" {
		t.Fatalf("unexpected key prefix: %s", l.KeyPrefix)
	}
}

func TestEnvironmentLoaderGetMap(t *testing.T) {
	const appFieldSetKey = "app"

	const idFieldKey = "id"

	const secretFieldKey = "secret"

	const appIDValue = "test-app-id"

	const appSecretValue = "test-sensitive-value"

	appIDEnvironmentVariable := strings.ToUpper(fmt.Sprintf("%s_%s", appFieldSetKey, idFieldKey))
	appSecretEnvironmentVariable := strings.ToUpper(fmt.Sprintf("%s_%s", appFieldSetKey, secretFieldKey))

	os.Setenv(appIDEnvironmentVariable, appIDValue)
	os.Setenv(appSecretEnvironmentVariable, appSecretValue)

	loader := bconf.NewEnvironmentLoader()

	appMap := loader.GetMap(appFieldSetKey, []string{idFieldKey, secretFieldKey, "invalid_field_key"})
	if len(appMap) != 2 {
		t.Fatalf("unexpected appMap length '%d', expected '2'", len(appMap))
	}
}

func TestEnvironmentLoader(t *testing.T) {
	const sessionFieldSet = "session"

	const sessionTokenKey = "token"

	const logFieldSet = "log"

	const logLevelKey = "level"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	envSessionTokenKey := strings.ToUpper(fmt.Sprintf("%s_%s", sessionFieldSet, sessionTokenKey))
	envLogLevelKey := strings.ToUpper(fmt.Sprintf("%s_%s", logFieldSet, logLevelKey))

	os.Setenv(envSessionTokenKey, sessionKeyValue)
	os.Setenv(envLogLevelKey, logLevelValue)

	l := bconf.EnvironmentLoader{}
	clone := l.Clone()

	if l.Name() != "bconf_environment" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionFieldSet, sessionTokenKey), envSessionTokenKey) {
		t.Errorf("unexpected loader help string contents: '%s'", l.HelpString(sessionFieldSet, sessionTokenKey))
	}

	if !strings.Contains(clone.HelpString(sessionFieldSet, sessionTokenKey), envSessionTokenKey) {
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

func TestEnvironmentLoaderWithKeyPrefix(t *testing.T) {
	const sessionFieldSet = "session"

	const sessionTokenKey = "token"

	const logFieldSet = "log"

	const logLevelKey = "level"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	const keyPrefix = "ext_http_api"

	envSessionKey := strings.ToUpper(fmt.Sprintf("%s_%s_%s", keyPrefix, sessionFieldSet, sessionTokenKey))
	envLogLevelKey := strings.ToUpper(fmt.Sprintf("%s_%s_%s", keyPrefix, logFieldSet, logLevelKey))

	os.Setenv(envSessionKey, sessionKeyValue)
	os.Setenv(envLogLevelKey, logLevelValue)

	l := bconf.EnvironmentLoader{
		KeyPrefix: keyPrefix,
	}
	clone := l.Clone()

	if l.Name() != "bconf_environment" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionFieldSet, sessionTokenKey), envSessionKey) {
		t.Errorf("unexpected loader help string contents: '%s'", clone.HelpString(sessionFieldSet, sessionTokenKey))
	}

	if !strings.Contains(clone.HelpString(sessionFieldSet, sessionTokenKey), envSessionKey) {
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
