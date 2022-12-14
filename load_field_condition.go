package bconf

import "fmt"

type FieldCondition struct {
	Condition   func(fieldValue any) bool
	FieldSetKey string
	FieldKey    string
}

func (c *FieldCondition) Clone() LoadCondition {
	clone := *c

	return &clone
}

func (c *FieldCondition) FieldDependency() (fieldSetKey, fieldKey string) {
	return c.FieldSetKey, c.FieldKey
}

func (c *FieldCondition) Load(value any) bool {
	return c.Condition(value)
}

func (c *FieldCondition) Validate() []error {
	errs := []error{}

	if c.FieldSetKey == "" {
		errs = append(errs, fmt.Errorf("field-set key required for field condition"))
	}

	if c.FieldKey == "" {
		errs = append(errs, fmt.Errorf("field key required for field condition"))
	}

	return errs
}
