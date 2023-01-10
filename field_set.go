package bconf

import "fmt"

type FieldSets []*FieldSet

type FieldSet struct {
	fieldMap       map[string]*Field
	Key            string
	LoadConditions LoadConditions
	Fields         Fields
}

func (f *FieldSet) Clone() *FieldSet {
	clone := *f

	if len(f.LoadConditions) > 0 {
		clone.LoadConditions = make([]LoadCondition, len(f.LoadConditions))
		for index, value := range f.LoadConditions {
			clone.LoadConditions[index] = value.Clone()
		}
	}

	if len(f.Fields) > 0 {
		clone.Fields = make([]*Field, len(f.Fields))

		for index, field := range f.Fields {
			newField := *field
			clone.Fields[index] = &newField
		}
	}

	if len(f.fieldMap) > 0 {
		clone.fieldMap = make(map[string]*Field, len(f.fieldMap))

		for key, field := range f.fieldMap {
			newField := field.Clone()
			clone.fieldMap[key] = newField
		}
	}

	return &clone
}

// validate validates the configuration of the field set.
func (f *FieldSet) validate() []error {
	errs := []error{}

	if f.Key == "" {
		errs = append(errs, fmt.Errorf("field-set key required"))
	}

	if len(f.LoadConditions) > 0 {
		for _, loadCondition := range f.LoadConditions {
			if loadConditionErrs := loadCondition.Validate(); len(loadConditionErrs) > 0 {
				for _, err := range loadConditionErrs {
					errs = append(errs, fmt.Errorf("load condition validation error: %w", err))
				}
			}
		}
	}

	fieldKeys := map[string]struct{}{}

	if len(f.Fields) > 0 {
		for _, field := range f.Fields {
			if _, found := fieldKeys[field.Key]; found {
				errs = append(errs, fmt.Errorf("duplicate field key found: '%s'", field.Key))
				continue
			}

			fieldKeys[field.Key] = struct{}{}
		}
	}

	return errs
}

// initializeFieldMap transitions the FieldSet over from its configuration []*Field to a map[string]*Field.
func (f *FieldSet) initializeFieldMap() {
	fieldMap := make(map[string]*Field, len(f.Fields))
	for _, field := range f.Fields {
		fieldMap[field.Key] = field
	}

	f.fieldMap = fieldMap
	f.Fields = nil
}

// generateFieldDefaults runs field default generators exactly once. Multiple calls will not regenerate field defaults.
func (f *FieldSet) generateFieldDefaults() []error {
	errs := []error{}

	if f.fieldMap != nil {
		for key, field := range f.fieldMap {
			if err := field.generateDefault(); err != nil {
				errs = append(errs, fmt.Errorf("field '%s' default value generation error: %w", key, err))
			}
		}
	}

	return errs
}

// validateFields validates field configuration, and can only be run after field defaults have been generated.
func (f *FieldSet) validateFields() []error {
	errs := []error{}

	if f.fieldMap != nil {
		for key, field := range f.fieldMap {
			if fieldErrs := field.validate(); len(fieldErrs) > 0 {
				for _, err := range fieldErrs {
					errs = append(errs, fmt.Errorf("field '%s' validation error: %w", key, err))
				}
			}
		}
	}

	return errs
}
