package bconf

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
	return nil
}
