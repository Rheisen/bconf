package bconf_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rheisen/bconf"
)

func TestFieldBuilderCreate(t *testing.T) {
	builder := bconf.FieldBuilder{}

	field := builder.Create()
	if field == nil {
		t.Fatalf("unexpected nil field from builder create")
	}

	field = bconf.NewFieldBuilder().Create()
	if field == nil {
		t.Fatalf("unexpected nil field from builder create")
	}

	field = bconf.FB().Create()
	if field == nil {
		t.Fatalf("unexpected nil field from builder create")
	}
}

func TestFieldBuilderKey(t *testing.T) {
	const fieldKey = "test_field_key"

	field := bconf.FB().Key(fieldKey).Create()
	if field.Key != fieldKey {
		t.Fatalf("unexpected field key '%s', expected '%s'", field.Key, fieldKey)
	}
}

func TestFieldBuilderDefault(t *testing.T) {
	const fieldDefault = 30 * time.Second

	field := bconf.FB().Default(fieldDefault).Create()
	if field.Default != fieldDefault {
		t.Fatalf("unexpected field default '%v', expected '%v'", field.Default, fieldDefault)
	}
}

func TestFieldBuilderValidator(t *testing.T) {
	validator := func(fieldValue any) error {
		return fmt.Errorf("validator error")
	}
	field := bconf.FB().Validator(validator).Create()

	if err := field.Validator(nil); err.Error() != "validator error" {
		t.Fatalf("unexpected validator error value: %s", err)
	}
}

func TestFieldBuilderDefaultGenerator(t *testing.T) {
	defaultGenerator := func() (any, error) {
		return "default", nil
	}
	field := bconf.FB().DefaultGenerator(defaultGenerator).Create()

	generatedDefault, _ := field.DefaultGenerator()
	if generatedDefault != "default" {
		t.Fatalf("unexpected generated default value '%s', expected 'default'", generatedDefault)
	}
}

func TestFieldBuilderType(t *testing.T) {
	fieldType := bconf.Float

	field := bconf.FB().Type(fieldType).Create()
	if field.Type != fieldType {
		t.Fatalf("unexpected field type '%s', expected '%s'", field.Type, fieldType)
	}
}

func TestFieldBuilderDescription(t *testing.T) {
	const fieldDescription = "field description test"

	field := bconf.FB().Description(fieldDescription).Create()
	if field.Description != fieldDescription {
		t.Fatalf("unexpected field description '%s', expected '%s'", field.Description, fieldDescription)
	}
}

func TestFieldBuilderEnumeration(t *testing.T) {
	fieldEnumeration := []any{"one", "two", "three"}
	field := bconf.FB().Enumeration(fieldEnumeration...).Create()

	if len(field.Enumeration) != len(fieldEnumeration) {
		t.Fatalf(
			"unexpected field enumeration length '%d', expected '%d'",
			len(field.Enumeration), len(fieldEnumeration),
		)
	}
}

func TestFieldBuilderRequired(t *testing.T) {
	field := bconf.FB().Required().Create()
	if field.Required == false {
		t.Fatalf("expected field to be required")
	}
}

func TestFieldBuilderSensitive(t *testing.T) {
	field := bconf.FB().Sensitive().Create()
	if field.Sensitive == false {
		t.Fatalf("expected field to be sensitive")
	}
}
