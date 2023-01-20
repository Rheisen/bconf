package bconf_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/rheisen/bconf"
)

func TestJSONFileLoaderFunctions(t *testing.T) {
	loader := bconf.NewJSONFileLoader()

	if loader == nil {
		t.Fatalf("unexpected nil loader")
	}

	loader = bconf.NewJSONFileLoaderWithAttributes(json.Unmarshal, "./fixtures/json_config_test_fixture_01.json")

	if loader == nil {
		t.Fatalf("unexpected nil loader")
	}

	if loader.Decoder == nil {
		t.Fatalf("unexpected nil decoder")
	}

	if len(loader.FilePaths) != 1 {
		t.Fatalf("unexpected file-paths length '%d', expected '1'", len(loader.FilePaths))
	}
}

func TestJSONFileLoaderClone(t *testing.T) {
	loader := loaderWithTestFixture01()
	clone := loader.Clone()

	if len(clone.FilePaths) != len(loader.FilePaths) {
		t.Fatalf("unexpected clone file-path length '%d', expected '%d'", len(clone.FilePaths), len(loader.FilePaths))
	}

	loader.FilePaths[0] = "./fixtures/json_config_test_fixture_02.json"

	if clone.FilePaths[0] == loader.FilePaths[0] {
		t.Fatalf("unexpected clone file-path value: %s", clone.FilePaths[0])
	}

	loader.FilePaths[0] = "./fixtures/json_config_test_fixture_01.json"

	loaderClone := loader.CloneLoader()

	loader.FilePaths[0] = "./fixtures/empty.json"

	_, found := loaderClone.Get("app", "id")
	if !found {
		t.Fatalf("unexpected issue finding app-id")
	}
}

func TestJSONFileLoaderName(t *testing.T) {
	loader := bconf.NewJSONFileLoader()

	if loader.Name() != "bconf_jsonfile" {
		t.Fatalf("unexpected json-file-loader name '%s'", loader.Name())
	}
}

func TestJSONFileLoaderGet(t *testing.T) {
	loaderFixture01 := loaderWithTestFixture01()
	loaderNoFilePaths := loaderWithNoFilePaths()
	loaderInvalidFilePaths := loaderWithInvalidFilePaths()
	loaderBadDecoder := loaderWithBadDecoder()

	_, found := loaderFixture01.Get("strange_key", "some_field")
	if found {
		t.Fatalf("unexpected found value when looking for non-existent key")
	}

	appID, found := loaderFixture01.Get("app", "id")

	if !found {
		t.Fatalf("expected loader with fixture file to find appID value")
	}

	if appID != "test-app-id" {
		t.Fatalf("unexpected appID value '%s', expected 'test-app-id'", appID)
	}

	appPort, found := loaderFixture01.Get("app", "port")
	if !found {
		t.Fatalf("expected loader with fixture file to find appPort value")
	}

	if appPort != "8080" {
		t.Fatalf("unexpected appPort value '%s', expected '8080'", appPort)
	}

	_, found = loaderNoFilePaths.Get("app", "id")
	if found {
		t.Fatalf("unexpected appID found by loader with no file-paths")
	}

	_, found = loaderInvalidFilePaths.Get("app", "id")
	if found {
		t.Fatalf("unexpected appID found by loader with invalid file-paths")
	}

	_, found = loaderBadDecoder.Get("app", "id")
	if found {
		t.Fatalf("unexpected appID found by loader with bad decoder")
	}
}

func TestJSONFileLoaderGetMap(t *testing.T) {
	loaderFixture01 := loaderWithTestFixture01()
	loaderNoFilePaths := loaderWithNoFilePaths()

	appMap := loaderFixture01.GetMap("app", []string{"id", "secret", "invalid_field_key"})
	if len(appMap) != 2 {
		t.Fatalf("unexpected length of app field-set map '%d', expected '2'", len(appMap))
	}

	appMap = loaderNoFilePaths.GetMap("app", []string{"id", "secret", "invalid_field_key"})
	if len(appMap) != 0 {
		t.Fatalf("unexpected length of app file-set map '%d', expected '0'", len(appMap))
	}
}

func TestJSONFileLoaderHelpString(t *testing.T) {
	loaderFixture01 := loaderWithTestFixture01()

	helpString := loaderFixture01.HelpString("app", "id")

	if !strings.Contains(helpString, "app.id") {
		t.Fatalf("unexpected help string: '%s'", helpString)
	}
}

func loaderWithTestFixture01() *bconf.JSONFileLoader {
	return bconf.NewJSONFileLoaderWithAttributes(json.Unmarshal, "./fixtures/json_config_test_fixture_01.json")
}

func loaderWithTestFixture02() *bconf.JSONFileLoader {
	return bconf.NewJSONFileLoaderWithAttributes(json.Unmarshal, "./fixtures/json_config_test_fixture_02.json")
}

func loaderWithBadDecoder() *bconf.JSONFileLoader {
	badDecoder := func(data []byte, v interface{}) error {
		return fmt.Errorf("decoder error")
	}

	return bconf.NewJSONFileLoaderWithAttributes(
		badDecoder, "./fixtures/json_config_test_fixture_01.json",
	)
}

func loaderWithNoFilePaths() *bconf.JSONFileLoader {
	return bconf.NewJSONFileLoader()
}

func loaderWithInvalidFilePaths() *bconf.JSONFileLoader {
	return bconf.NewJSONFileLoaderWithAttributes(nil, "./fixtures/non-existent-file.json")
}
