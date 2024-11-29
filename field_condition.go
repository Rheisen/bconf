package bconf

import "fmt"

type FieldCondition struct {
	Condition   func(fieldValues map[string]any) (bool, error)
	FieldSetKey string
	FieldKeys   []string
}

func (c *FieldCondition) Clone() LoadCondition {
	clone := *c

	return &clone
}

func (c *FieldCondition) FieldDependency() (fieldSetKey string, fieldKeys []string) {
	return c.FieldSetKey, c.FieldKeys
}

func (c *FieldCondition) Load(values map[string]any) (bool, error) {
	return c.Condition(values)
}

func (c *FieldCondition) Validate() []error {
	errs := []error{}

	if c.FieldSetKey == "" {
		errs = append(errs, fmt.Errorf("field-set key required for field condition"))
	}

	if len(c.FieldKeys) == 0 {
		errs = append(errs, fmt.Errorf("at least one field key required for field condition"))
	}

	return errs
}
