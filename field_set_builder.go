package bconf

func NewFieldSetBuilder() *FieldSetBuilder {
	return &FieldSetBuilder{fieldSet: &FieldSet{}}
}

func FSB() *FieldSetBuilder {
	return NewFieldSetBuilder()
}

type FieldSetBuilder struct {
	fieldSet *FieldSet
}

func (b *FieldSetBuilder) Key(value string) *FieldSetBuilder {
	b.init()
	b.fieldSet.Key = value

	return b
}

func (b *FieldSetBuilder) Fields(value ...*Field) *FieldSetBuilder {
	b.init()
	b.fieldSet.Fields = value

	return b
}

func (b *FieldSetBuilder) LoadConditions(value ...LoadCondition) *FieldSetBuilder {
	b.init()
	b.fieldSet.LoadConditions = value

	return b
}

func (b *FieldSetBuilder) Create() *FieldSet {
	b.init()
	return b.fieldSet.Clone()
}

func (b *FieldSetBuilder) init() {
	if b.fieldSet == nil {
		b.fieldSet = &FieldSet{}
	}
}
