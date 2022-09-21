package bconf_test

import (
	"strings"
	"testing"

	"github.com/rheisen/bconf"
	"github.com/rheisen/bconf/bconfconst"
)

func TestStringField(t *testing.T) {
	validField := bconf.Field{
		FieldType:   bconfconst.String,
		Required:    true,
		Default:     "value",
		Help:        "basic field",
		Enumeration: []any{"arg1", "arg2", "value"},
	}

	errs := validField.Validate()
	if errs != nil && len(errs) > 0 {
		t.Errorf("unexpected errors validating validField: %v", errs)
	}

	invalidFieldInvalidDefaultType := bconf.Field{
		FieldType:   bconfconst.String,
		Required:    false,
		Default:     2,
		Help:        "basic field",
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
		Help:        "basic field",
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
