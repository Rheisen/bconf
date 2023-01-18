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

func TestEnvironmentLoader(t *testing.T) {
	const sessionKey = "session_key"

	const logLevelKey = "log_level"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	envSessionKey := strings.ToUpper(sessionKey)
	envLogLevelKey := strings.ToUpper(logLevelKey)

	os.Setenv(envSessionKey, sessionKeyValue)
	os.Setenv(envLogLevelKey, logLevelValue)

	l := bconf.EnvironmentLoader{}
	clone := l.Clone()

	if l.Name() != "bconf_environment" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionKey), envSessionKey) {
		t.Errorf("unexpected loader help string contents: '%s'", clone.HelpString(sessionKey))
	}

	if !strings.Contains(clone.HelpString(sessionKey), envSessionKey) {
		t.Errorf("unexpected loader clone help string contents: '%s'", clone.HelpString(sessionKey))
	}

	sessionKeyLookup, found := l.Get(sessionKey)
	if !found {
		t.Errorf("unexpected problem getting session_key value")
	}

	if sessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session key: '%s'", sessionKeyLookup)
	}

	logLevel, found := l.Get(logLevelKey)
	if !found {
		t.Errorf("unexpected problem getting log_level value")
	}

	if logLevel != logLevelValue {
		t.Errorf("unexpected value for log level: '%s'", logLevel)
	}

	cloneSessionKeyLookup, cloneFound := clone.Get(sessionKey)
	if !cloneFound {
		t.Errorf("unexpected problem getting session_key value from loader clone")
	}

	if cloneSessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session_key from loader clone: '%s'", cloneSessionKeyLookup)
	}
}

func TestEnvironmentLoaderWithKeyPrefix(t *testing.T) {
	const sessionKey = "session_key"

	const logLevelKey = "log_level"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	const keyPrefix = "ext_http_api"

	envSessionKey := strings.ToUpper(fmt.Sprintf("%s_%s", keyPrefix, sessionKey))
	envLogLevelKey := strings.ToUpper(fmt.Sprintf("%s_%s", keyPrefix, logLevelKey))

	os.Setenv(envSessionKey, sessionKeyValue)
	os.Setenv(envLogLevelKey, logLevelValue)

	l := bconf.EnvironmentLoader{
		KeyPrefix: keyPrefix,
	}
	clone := l.Clone()

	if l.Name() != "bconf_environment" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionKey), envSessionKey) {
		t.Errorf("unexpected loader help string contents: '%s'", clone.HelpString(sessionKey))
	}

	if !strings.Contains(clone.HelpString(sessionKey), envSessionKey) {
		t.Errorf("unexpected loader clone help string contents: '%s'", clone.HelpString(sessionKey))
	}

	sessionKeyLookup, found := l.Get(sessionKey)
	if !found {
		t.Errorf("unexpected problem getting session_key value")
	}

	if sessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session key: '%s'", sessionKeyLookup)
	}

	logLevel, found := l.Get(logLevelKey)
	if !found {
		t.Errorf("unexpected problem getting log_level value")
	}

	if logLevel != logLevelValue {
		t.Errorf("unexpected value for log level: '%s'", logLevel)
	}

	cloneSessionKeyLookup, cloneFound := clone.Get(sessionKey)
	if !cloneFound {
		t.Errorf("unexpected problem getting session_key value from loader clone")
	}

	if cloneSessionKeyLookup != sessionKeyValue {
		t.Errorf("unexpected value for session_key from loader clone: '%s'", cloneSessionKeyLookup)
	}
}
