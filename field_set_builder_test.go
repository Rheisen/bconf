package bconf_test

import (
	"testing"

	"github.com/rheisen/bconf"
)

func TestFieldSetBuilderCreate(t *testing.T) {
	fieldSet := bconf.NewFieldSetBuilder().Create()
	if fieldSet == nil {
		t.Fatalf("unexpected nil field-set")
	}

	fieldSet = bconf.FSB().Create()
	if fieldSet == nil {
		t.Fatalf("unexpected nil field-set")
	}

	builder := &bconf.FieldSetBuilder{}

	fieldSet = builder.Create()
	if fieldSet == nil {
		t.Fatalf("unexpected nil field-set")
	}
}

func TestFieldSetBuilderKey(t *testing.T) {
	fieldSetKey := "test_key"

	fieldSet := bconf.FSB().Key(fieldSetKey).Create()
	if fieldSet.Key != fieldSetKey {
		t.Fatalf("unexpected field-set key '%s', expected '%s'", fieldSet.Key, fieldSetKey)
	}
}

func TestFieldSetBuilderFields(t *testing.T) {
	fieldKey := "test_field_key"
	field := bconf.FB().Key(fieldKey).Create()

	fieldSet := bconf.FSB().Fields(field).Create()
	if len(fieldSet.Fields) != 1 {
		t.Fatalf("unexpected fields length '%d', expected 1", len(fieldSet.Fields))
	}

	if fieldSet.Fields[0].Key != fieldKey {
		t.Fatalf("unexpected field key '%s', expected '%s'", fieldSet.Fields[0].Key, fieldKey)
	}
}

func TestFieldSetBuilderLoadConditions(t *testing.T) {
	loadConditionFieldSetKey := "test_field_set_key"
	loadConditionFieldKey := "test_field_key"

	fieldSet := bconf.FSB().LoadConditions(
		bconf.FCB().FieldSetKey(loadConditionFieldSetKey).FieldKey(loadConditionFieldKey).Create(),
	).Create()

	if len(fieldSet.LoadConditions) != 1 {
		t.Fatalf("unexpected load-conditions length '%d', expected 1", len(fieldSet.LoadConditions))
	}

	fieldSetKey, fieldKey := fieldSet.LoadConditions[0].FieldDependency()
	if fieldSetKey != loadConditionFieldSetKey {
		t.Fatalf("unexpected field-set key '%s', expected '%s'", fieldSetKey, loadConditionFieldSetKey)
	}

	if fieldKey != loadConditionFieldKey {
		t.Fatalf("unexpected field key '%s', expected '%s'", fieldKey, loadConditionFieldKey)
	}
}
