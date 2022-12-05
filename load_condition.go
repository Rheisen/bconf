package bconf

type LoadCondition interface {
	Clone() LoadCondition
	FieldDependency() (string, string)
	Load(value any) bool
	Validate() []error
}
