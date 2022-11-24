package bconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rheisen/bconf/bconfconst"
)

type Field struct {
	FieldType        string
	Required         bool
	Default          any
	Description      string
	Enumeration      []any
	Validator        func(v any) error
	DefaultGenerator func() (any, error)
	// -- private attributes --
	fieldFound       []string
	fieldValue       map[string]any
	generatedDefault any
}

func (f *Field) GenerateDefault() error {
	if f.DefaultGenerator == nil {
		return nil
	}

	generatedDefault, err := f.DefaultGenerator()
	if err != nil {
		return fmt.Errorf("problem generating default field value: %w", err)
	}

	f.generatedDefault = generatedDefault

	return nil
}

func (f *Field) Validate() []error {
	errs := []error{}

	fieldTypeFound := false
	for _, fieldType := range bconfconst.FieldTypes() {
		if fieldType == f.FieldType {
			fieldTypeFound = true

			// Check that default and default generator are not both set
			if f.Default != nil && f.DefaultGenerator != nil {
				errs = append(errs, fmt.Errorf(bconfconst.ErrorFieldDefaultSetting))
			}

			// Check that required and default are not both set
			if f.Required && f.Default != nil || f.Required && f.DefaultGenerator != nil {
				errs = append(errs, fmt.Errorf(bconfconst.ErrorFieldRequiredWithDefault))
			}

			// Check that the type of the default value matches the field type
			if f.Default != nil {
				if reflect.TypeOf(f.Default).String() != fieldType {
					errs = append(
						errs,
						fmt.Errorf(
							"invalid default type: expected '%s', found '%s'",
							fieldType,
							reflect.TypeOf(f.Default).String(),
						),
					)
				}
			}
			if f.generatedDefault != nil {
				if reflect.TypeOf(f.generatedDefault).String() != fieldType {
					errs = append(
						errs,
						fmt.Errorf(
							"invalid generated default type: expected '%s', found '%s'",
							fieldType,
							reflect.TypeOf(f.generatedDefault).String(),
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
								"invalid enumeration value type: expected '%s', found '%s'",
								fieldType,
								reflect.TypeOf(val).String(),
							),
						)
					}
				}
			}

			// Return here before validating default values existing in enumeration list
			if len(errs) > 0 {
				return errs
			}

			// Check that a given default value is in the list of enumerated acceptable values
			if f.Default != nil && !f.valueInEnumeration(f.Default) {
				errs = append(
					errs,
					fmt.Errorf(
						"invalid default value: default value '%s' expected in enumeration list",
						f.Default,
					),
				)
			}
			if f.generatedDefault != nil && !f.valueInEnumeration(f.generatedDefault) {
				errs = append(
					errs,
					fmt.Errorf(
						"invalid generated default value: default value '%s' expected in enumeration list",
						f.Default,
					),
				)
			}
		}
	}

	if !fieldTypeFound {
		errs = append(errs, fmt.Errorf("invalid field type specified: '%s'", f.FieldType))
	}

	return errs
}

func (f *Field) GetValue() (any, error) {
	if f.fieldFound != nil {
		value, found := f.fieldValue[f.fieldFound[len(f.fieldFound)-1]]
		if !found {
			return nil, fmt.Errorf("library error, please report")
		}

		return value, nil
	}
	if f.Default != nil {
		return f.Default, nil
	}
	if f.generatedDefault != nil {
		return f.generatedDefault, nil
	}

	return nil, fmt.Errorf("empty field value")
}

func (f *Field) GetValueFrom(loader string) (any, error) {
	if f.fieldValue == nil {
		return nil, fmt.Errorf("")
	}

	value, found := f.fieldValue[loader]
	if !found {
		return nil, fmt.Errorf("")
	}

	return value, nil
}

func (f *Field) set(loaderName string, value string) error {
	parsedValue, err := f.parseString(value)
	if err != nil {
		return fmt.Errorf("problem parsing field value to field type: %w", err)
	}

	if !f.valueInEnumeration(parsedValue) {
		return fmt.Errorf("parsed value not found in enumeration list")
	}

	if f.fieldFound == nil {
		f.fieldFound = []string{loaderName}
	} else {
		f.fieldFound = append(f.fieldFound, loaderName)
	}

	if f.fieldValue == nil {
		f.fieldValue = map[string]any{loaderName: parsedValue}
	} else {
		f.fieldValue[loaderName] = value
	}

	return nil
}

func (f *Field) parseString(value string) (any, error) {
	switch f.FieldType {
	case bconfconst.String:
		return value, nil
	case bconfconst.Bool:
		return strconv.ParseBool(value)
	case bconfconst.Int:
		return strconv.Atoi(value)
	case bconfconst.Time:
		return time.Parse(time.RFC3339, value)
	case bconfconst.Duration:
		return time.ParseDuration(value)
	default:
		return "", fmt.Errorf("unsupported field type: %s", f.FieldType)
	}
}

func (f *Field) valueInEnumeration(value any) bool {
	if len(f.Enumeration) < 1 {
		return true
	}

	for _, acceptedValue := range f.Enumeration {
		if value == acceptedValue {
			return true
		}
	}

	return false
}

func (f *Field) enumerationString() string {
	builder := strings.Builder{}

	if len(f.Enumeration) > 0 {
		builder.WriteString("[")
		for index, value := range f.Enumeration {
			if index != 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("'%s'", value))
		}
		builder.WriteString("]")
	}

	return builder.String()
}
