package bconf_test

import (
	"testing"

	"github.com/rheisen/bconf"
)

func TestFieldConditionBuilderCreate(t *testing.T) {
	builder := &bconf.FieldConditionBuilder{}

	condition := builder.Create()
	if condition == nil {
		t.Fatalf("unexpected nil condition")
	}

	condition = bconf.NewFieldConditionBuilder().Create()
	if condition == nil {
		t.Fatalf("unexpected nil condition")
	}

	condition = bconf.FCB().Create()
	if condition == nil {
		t.Fatalf("unexpected nil condition")
	}
}

func TestFieldConditionBuilderKeys(t *testing.T) {
	const fieldSetKey = "test_field_set_key"

	const fieldKey = "test_field_key"

	condition := bconf.FCB().FieldSetKey(fieldSetKey).Create()
	if fsKey, _ := condition.FieldDependency(); fsKey != fieldSetKey {
		t.Fatalf("unexpected field-set key value '%s', expected '%s'", fsKey, fieldSetKey)
	}

	condition = bconf.FCB().FieldKey(fieldKey).Create()
	if _, fKey := condition.FieldDependency(); fKey != fieldKey {
		t.Fatalf("unexpected field key value '%s', expected '%s'", fKey, fieldKey)
	}
}

func TestFieldConditionBuilderCondition(t *testing.T) {
	condition := func(fieldValue any) (bool, error) {
		return true, nil
	}
	fieldCondition := bconf.FCB().Condition(condition).Create()

	if ok, _ := fieldCondition.Load(nil); !ok {
		t.Fatalf("unexpected load value: %v", ok)
	}
}
