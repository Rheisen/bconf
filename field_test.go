package bconf_test

import (
	"strings"
	"testing"
	"time"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

func TestFieldTypes(t *testing.T) {
	type FieldTypeTest struct {
		FieldType       string
		ValidDefaults   []any
		InvalidDefaults []any
	}

	testCases := map[string]FieldTypeTest{
		"bool-field-with-default": {
			FieldType:       bconfconst.Bool,
			ValidDefaults:   []any{true, false},
			InvalidDefaults: []any{},
		},
		"string-field-with-default": {
			FieldType:     bconfconst.String,
			ValidDefaults: []any{"string-default", "", "-"},
		},
		"int-field-with-default": {
			FieldType:     bconfconst.Int,
			ValidDefaults: []any{-512, -256, 0, 256, 512},
		},
		"float-64-field-with-default": {
			FieldType:       bconfconst.Float64,
			ValidDefaults:   []any{-1024.64, -512.0, 0.0, 512.0, 1024.64},
			InvalidDefaults: []any{1024, 0, "1024", true, false},
		},
		"time-field-with-default": {
			FieldType:     bconfconst.Time,
			ValidDefaults: []any{time.Now()},
		},
	}

	for name, testCase := range testCases {
		for _, expectedValidDefault := range testCase.ValidDefaults {
			field := bconf.Field{
				FieldType: testCase.FieldType,
				Default:   expectedValidDefault,
			}

			if errs := field.Validate(); len(errs) > 0 {
				t.Errorf("TestFieldTypes test case %s: validation error(s): %v", name, errs)
			}

			if fieldValue, err := field.GetValue(); err != nil {
				t.Errorf("TestFieldTypes test case %s: get value error: %s", name, err)
			} else if fieldValue != expectedValidDefault {
				t.Errorf("TestFieldTypes test case %s: unexpected value: %s", name, fieldValue)
			}
		}

		for _, expectedInvalidDefault := range testCase.InvalidDefaults {
			field := bconf.Field{
				FieldType: testCase.FieldType,
				Default:   expectedInvalidDefault,
			}

			if errs := field.Validate(); len(errs) < 1 {
				t.Errorf("TestFieldTypes test case %s: expected validation errors", name)
			}
		}
	}
}

func TestStringField(t *testing.T) {
	validField := bconf.Field{
		FieldType:   bconfconst.String,
		Default:     "value",
		Description: "basic field",
		Enumeration: []any{"arg1", "arg2", "value"},
	}

	errs := validField.Validate()
	if len(errs) > 0 {
		t.Errorf("unexpected errors validating validField: %v", errs)
	}

	invalidFieldInvalidDefaultType := bconf.Field{
		FieldType:   bconfconst.String,
		Default:     2,
		Description: "basic field",
		Enumeration: []any{"arg1", "arg2", "value"},
	}

	errs = invalidFieldInvalidDefaultType.Validate()
	if errs == nil || len(errs) == 0 {
		t.Errorf("expected 'invalid default type' error for invalidFieldInvalidDefaultType, found no errors")
	}
	if errs != nil && len(errs) > 0 && !strings.Contains(errs[0].Error(), "invalid default type") {
		t.Errorf("expected 'invalid default type' error, found '%s'", errs[0].Error())
	}
	if errs != nil && len(errs) > 1 {
		t.Errorf("expected 1 err for invalidFieldInvalidDefaultType, found %d (%v)", len(errs), errs)
	}

	invalidFieldInvalidEnumerationList := bconf.Field{
		FieldType:   bconfconst.String,
		Required:    false,
		Default:     "value",
		Description: "basic field",
		Enumeration: []any{1, 2, "3", 32.812, "forty-five", "value", true},
	}

	errs = invalidFieldInvalidEnumerationList.Validate()
	if errs == nil || len(errs) == 0 {
		t.Errorf(
			"expected 'invalid enumeration value type' error for invalidFieldInvalidEnumerationList, found no errors",
		)
	}
	if errs != nil && len(errs) != 4 {
		t.Errorf("expected 4 errs for invalidFieldInvalidEnumerationList, found %d (%v)", len(errs), errs)
	}
	if errs != nil && len(errs) > 0 {
		for _, err := range errs {
			if !strings.Contains(err.Error(), "invalid enumeration value type") {
				t.Errorf("expected 'invalid enumeration value type' error, found '%s'", err.Error())
			}
		}
	}
}

func TestFieldDefaultGenerator(t *testing.T) {
	validField := bconf.Field{
		FieldType: bconfconst.String,
		DefaultGenerator: func() (any, error) {
			return "generated-default", nil
		},
		Enumeration: []any{"generated-default", "other", "and-another"},
	}

	if err := validField.GenerateDefault(); err != nil {
		t.Fatalf("unexpected error generating default: %s", err)
	}

	if errs := validField.Validate(); len(errs) > 0 {
		t.Fatalf("unexpected errors validating validField: %v", errs)
	}

	if val, err := validField.GetValue(); err != nil {
		t.Fatalf("unexpected errors getting field value: %s", err)
	} else {
		t.Log(val)
	}
}
