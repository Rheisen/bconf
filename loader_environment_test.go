package bconf_test

import (
	"os"
	"strings"
	"testing"

	"github.com/rheisen/bconf"
)

func TestEnvironmentLoader(t *testing.T) {
	l := bconf.EnvironmentLoader{
		KeyPrefix: "ext_http_api",
	}
	clone := l.Clone()

	if l.Name() != "bconf_environment" {
		t.Errorf("unexpected loader name: '%s'", l.Name())
	}

	if !strings.Contains(l.HelpString("session_key"), strings.ToUpper("ext_http_api_session_key")) {
		t.Errorf("unexpected loader help string contents: '%s'", clone.HelpString("session_key"))
	}

	if !strings.Contains(clone.HelpString("session_key"), strings.ToUpper("ext_http_api_session_key")) {
		t.Errorf("unexpected loader clone help string contents: '%s'", clone.HelpString("session_key"))
	}

	os.Setenv("EXT_HTTP_API_SESSION_KEY", "abc123")
	os.Setenv("EXT_HTTP_API_LOG_LEVEL", "error")

	sessionKey, found := l.Get("session_key")
	if !found {
		t.Errorf("unexpected problem getting session_key value")
	}

	if sessionKey != "abc123" {
		t.Errorf("unexpected value for session key: '%s'", sessionKey)
	}

	logLevel, found := l.Get("log_level")
	if !found {
		t.Errorf("unexpected problem getting log_level value")
	}

	if logLevel != "error" {
		t.Errorf("unexpected value for log level: '%s'", logLevel)
	}

	cloneSessionKey, cloneFound := l.Get("session_key")
	if !cloneFound {
		t.Errorf("unexpected problem getting session_key value from loader clone")
	}

	if cloneSessionKey != "abc123" {
		t.Errorf("unexpected value for session_key from loader clone: '%s'", cloneSessionKey)
	}
}
