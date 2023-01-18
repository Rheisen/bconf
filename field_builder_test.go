package bconf_test

import (
	"testing"

	"github.com/rheisen/bconf"
)

func TestFieldBuilderInitialization(t *testing.T) {
	builder := bconf.FieldBuilder{}

	field := builder.Create()
	if field == nil {
		t.Errorf("unexpected nil field from builder create")
	}

	field = builder.Key("field_key").Create()
	if field.Key != "field_key" {
		t.Errorf("unexpected field.Key: %s", field.Key)
	}

	if field.Required == true {
		t.Errorf("unexpected field.Required value: '%v', expected 'true'", field.Required)
	}

	// fieldSets := bconf.FieldSets{
	// 	bconf.FSB().Key("api_config").Fields(
	// 		bconf.FB().Key("read_timeout").Type(bconf.Duration).Default(30*time.Second).Create(),
	// 		bconf.FB().Key("write_timeout").Type(bconf.Duration).Default(30*time.Second).Create(),
	// 		bconf.FB().Key("port").Type(bconf.Int).Default(8080).Create(),
	// 	).Create(),
	// }

	// if fieldSets == nil {
	// 	t.Errorf("")
	// }

	field = bconf.NewFieldBuilder().Key("field_key").Required().Create()

	if field.Required != true {
		t.Errorf("unexpected field.Required value: '%v', expected 'true'", field.Required)
	}

	field = bconf.FB().Enumeration("debug", "info", "warn", "error").Create()
}
