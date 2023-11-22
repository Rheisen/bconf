package bconf

type LoadConditions []LoadCondition

type LoadCondition interface {
	Clone() LoadCondition
	FieldDependency() (fieldSetKey string, fieldKeys []string)
	Load(values map[string]any) (bool, error)
	Validate() []error
}
