package bconf

type FieldCondition struct {
	FieldSetKey string
	FieldKey    string
	Condition   func(fieldValue any) bool
}

func (c *FieldCondition) Clone() LoadCondition {
	clone := *c

	return &clone
}

func (c *FieldCondition) FieldDependency() (string, string) {
	return c.FieldSetKey, c.FieldKey
}

func (c *FieldCondition) Load(value any) bool {
	return c.Condition(value)
}

func (c *FieldCondition) Validate() []error {
	return nil
}
