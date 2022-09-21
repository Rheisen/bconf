package bconf

import (
	"fmt"
	"reflect"

	"github.com/rheisen/bconf/bconfconst"
)

type Field struct {
	FieldType        string
	Required         bool
	Default          any
	Help             string
	Enumeration      []any
	Validator        func(v any) error
	DefaultGenerator func() (any, error)
	// -- private attributes --
	fieldFound []string
	fieldValue map[string]any
}

func (f *Field) Validate() []error {
	errs := []error{}

	fieldTypeFound := false
	for _, fieldType := range bconfconst.FieldTypes() {
		if fieldType == f.FieldType {
			fieldTypeFound = true
			// Check that the type of the default value matches the field type
			if f.Default != nil {
				if reflect.TypeOf(f.Default).String() != fieldType {
					errs = append(
						errs,
						fmt.Errorf(
							"invalid default type -- expected '%s', found '%s'",
							fieldType,
							reflect.TypeOf(f.Default).String(),
						),
					)
				}
			}
			// Check that the type of the enumeration list matches the field type
			if f.Enumeration != nil && len(f.Enumeration) > 0 {
				for _, val := range f.Enumeration {
					if reflect.TypeOf(val).String() != fieldType {
						errs = append(
							errs,
							fmt.Errorf(
								"invalid enumeration value type -- expected '%s', found '%s'",
								fieldType,
								reflect.TypeOf(val).String(),
							),
						)
					}
				}
			}
			// Check that a given default value is in the list of enumerated acceptable values
			if len(errs) == 0 && f.Enumeration != nil && len(f.Enumeration) > 0 && f.Default != nil {
				foundDefault := false
				for _, value := range f.Enumeration {
					if value == f.Default {
						foundDefault = true
					}
				}

				if !foundDefault {
					errs = append(
						errs,
						fmt.Errorf(
							"invalid default value -- default value '%s' expected in enumeration list",
							f.Default,
						),
					)
				}
			}
		}
	}

	if !fieldTypeFound {
		errs = append(errs, fmt.Errorf("invalid field type specified -- '%s'", f.FieldType))
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func GetValue() (any, error) {
	return nil, nil
}

func GetValueFrom(configLoader string) (any, error) {
	return nil, nil
}
