package bconf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rheisen/bconf"
)

func TestFlagLoader(t *testing.T) {
	const sessionKey = "session_key"

	const logLevelKey = "log_level"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	l := bconf.FlagLoader{
		OverrideLookup: []string{
			fmt.Sprintf("--%s=%s", sessionKey, sessionKeyValue),
			fmt.Sprintf("--%s", logLevelKey),
			logLevelValue,
		},
	}
	clone := l.Clone()

	if l.Name() != "bconf_flags" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionKey), sessionKey) {
		t.Errorf("unexpected loader help string contents: '%s'", clone.HelpString(sessionKey))
	}

	if !strings.Contains(clone.HelpString(sessionKey), sessionKey) {
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

func TestFlagLoaderWithKeyPrefix(t *testing.T) {
	const sessionKey = "session_key"

	const logLevelKey = "log_level"

	const sessionKeyValue = "abc123"

	const logLevelValue = "error"

	const keyPrefix = "ext_http_api"

	l := bconf.FlagLoader{
		KeyPrefix: keyPrefix,
		OverrideLookup: []string{
			fmt.Sprintf("--%s=%s", sessionKey, sessionKeyValue),
			fmt.Sprintf("--%s", logLevelKey),
			logLevelValue,
		},
	}
	clone := l.Clone()

	if l.Name() != "bconf_flags" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString(sessionKey), fmt.Sprintf("%s_%s", keyPrefix, sessionKey)) {
		t.Errorf("unexpected loader help string contents: '%s'", clone.HelpString(sessionKey))
	}

	if !strings.Contains(clone.HelpString(sessionKey), fmt.Sprintf("%s_%s", keyPrefix, sessionKey)) {
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
