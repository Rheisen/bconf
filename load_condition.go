package bconf

type LoadCondition interface {
	Clone() LoadCondition
	FieldDependency() (fieldSetKey string, fieldKey string)
	Load(value any) bool
	Validate() []error
}
