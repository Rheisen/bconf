package bconf

func FB() *FieldBuilder {
	return NewFieldBuilder()
}

func NewFieldBuilder() *FieldBuilder {
	return &FieldBuilder{field: &Field{}}
}

type FieldBuilder struct {
	field *Field
}

func (b *FieldBuilder) Key(value string) *FieldBuilder {
	b.init()
	b.field.Key = value

	return b
}

func (b *FieldBuilder) Default(value any) *FieldBuilder {
	b.init()
	b.field.Default = value

	return b
}

func (b *FieldBuilder) Validator(value func(fieldValue any) error) *FieldBuilder {
	b.init()
	b.field.Validator = value

	return b
}

func (b *FieldBuilder) DefaultGenerator(value func() (any, error)) *FieldBuilder {
	b.init()
	b.field.DefaultGenerator = value

	return b
}

func (b *FieldBuilder) Type(value string) *FieldBuilder {
	b.init()
	b.field.Type = value

	return b
}

func (b *FieldBuilder) Description(value string) *FieldBuilder {
	b.init()
	b.field.Description = value

	return b
}

func (b *FieldBuilder) Enumeration(value ...any) *FieldBuilder {
	b.init()
	b.field.Enumeration = value

	return b
}

func (b *FieldBuilder) Required() *FieldBuilder {
	b.init()
	b.field.Required = true

	return b
}

func (b *FieldBuilder) Sensitive() *FieldBuilder {
	b.init()
	b.field.Sensitive = true

	return b
}

func (b *FieldBuilder) Create() *Field {
	b.init()
	return b.field.Clone()
}

func (b *FieldBuilder) init() {
	if b.field == nil {
		b.field = &Field{}
	}
}
