package bconf

func FCB() *FieldConditionBuilder {
	return NewFieldConditionBuilder()
}

func NewFieldConditionBuilder() *FieldConditionBuilder {
	return &FieldConditionBuilder{condition: &FieldCondition{}}
}

type FieldConditionBuilder struct {
	condition *FieldCondition
}

func (b *FieldConditionBuilder) FieldSetKey(value string) *FieldConditionBuilder {
	b.init()
	b.condition.FieldSetKey = value

	return b
}

func (b *FieldConditionBuilder) AddFieldKey(value string) *FieldConditionBuilder {
	b.init()
	b.condition.FieldKeys = append(b.condition.FieldKeys, value)

	return b
}

func (b *FieldConditionBuilder) Condition(value func(fieldValues map[string]any) (bool, error)) *FieldConditionBuilder {
	b.init()
	b.condition.Condition = value

	return b
}

func (b *FieldConditionBuilder) Create() LoadCondition {
	b.init()
	return b.condition.Clone()
}

func (b *FieldConditionBuilder) init() {
	if b.condition == nil {
		b.condition = &FieldCondition{}
	}
}
