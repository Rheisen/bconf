package bconf

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewAppConfig(appName, appDescription string) *AppConfig {
	return &AppConfig{
		appName:          appName,
		appDescription:   appDescription,
		fieldSets:        map[string]*FieldSet{},
		orderedFieldSets: FieldSets{},
		loaders:          []Loader{},
	}
}

type AppConfig struct {
	fieldSets        map[string]*FieldSet
	appName          string
	appDescription   string
	loaders          []Loader
	orderedFieldSets FieldSets
	fieldSetLock     sync.Mutex
	register         sync.Once
	registered       bool
}

func (c *AppConfig) AppName() string {
	return c.appName
}

func (c *AppConfig) AppDescription() string {
	return c.appDescription
}

func (c *AppConfig) SetLoaders(loaders ...Loader) []error {
	errs := []error{}

	clonedLoaders := make([]Loader, len(loaders))
	for index, loader := range loaders {
		clonedLoaders[index] = loader.CloneLoader()
	}

	loaderNames := make(map[string]struct{}, len(clonedLoaders))
	for _, loader := range clonedLoaders {
		if _, found := loaderNames[loader.Name()]; found {
			errs = append(errs, fmt.Errorf("duplicate loader name found: '%s'", loader.Name()))
		}

		loaderNames[loader.Name()] = struct{}{}
	}

	if len(errs) > 0 {
		return errs
	}

	c.loaders = clonedLoaders

	return nil
}

func (c *AppConfig) AddFieldSet(fieldSet *FieldSet) []error {
	return c.addFieldSet(fieldSet, true)
}

func (c *AppConfig) AddFieldSets(fieldSets ...*FieldSet) []error {
	c.fieldSetLock.Lock()
	defer c.fieldSetLock.Unlock()

	errs := []error{}
	addedFieldSets := []string{}

	for _, fieldSet := range fieldSets {
		if fieldSetErrs := c.addFieldSet(fieldSet, false); len(fieldSetErrs) > 0 {
			errs = append(errs, fieldSetErrs...)
			continue
		}

		addedFieldSets = append(addedFieldSets, fieldSet.Key)
	}

	if len(errs) > 0 {
		for _, fieldSetKey := range addedFieldSets {
			delete(c.fieldSets, fieldSetKey)
		}

		c.orderedFieldSets = c.orderedFieldSets[:len(c.orderedFieldSets)-len(addedFieldSets)]
	}

	return errs
}

func (c *AppConfig) AddField(fieldSetKey string, field *Field) []error {
	c.fieldSetLock.Lock()
	defer c.fieldSetLock.Unlock()

	fieldSet, fieldSetFound := c.fieldSets[fieldSetKey]

	if !fieldSetFound {
		return []error{fmt.Errorf("no field-set found with key '%s'", fieldSetKey)}
	}

	if _, keyFound := c.fieldSets[fieldSetKey].fieldMap[field.Key]; keyFound {
		return []error{fmt.Errorf("duplicate field key found: '%s'", field.Key)}
	}

	field = field.Clone()

	if err := field.generateDefault(); err != nil {
		return []error{fmt.Errorf("field default value generation error: %w", err)}
	}

	if validationErrors := field.validate(); len(validationErrors) > 0 {
		return validationErrors
	}

	if err := c.checkForFieldDependencies(field, fieldSet); err != nil {
		return []error{fmt.Errorf("field dependency error: %w", err)}
	}

	c.fieldSets[fieldSetKey].fieldMap[field.Key] = field

	return nil
}

func (c *AppConfig) LoadFieldSet(fieldSetKey string) []error {
	errs := []error{}

	if !c.registered {
		errs = append(errs, fmt.Errorf("LoadFieldSet cannot be called before the app-config has been registered"))
		return errs
	}

	return c.loadFieldSet(fieldSetKey)
}

func (c *AppConfig) LoadField(fieldSetKey, fieldKey string) []error {
	errs := []error{}

	if !c.registered {
		errs = append(errs, fmt.Errorf("LoadField cannot be called before the app-config has been registered"))
		return errs
	}

	if _, fieldSetFound := c.fieldSets[fieldSetKey]; !fieldSetFound {
		errs = append(errs, fmt.Errorf("field-set with key '%s' not found", fieldSetKey))
		return errs
	}

	field, fieldKeyFound := c.fieldSets[fieldSetKey].fieldMap[fieldKey]
	if !fieldKeyFound {
		errs = append(errs, fmt.Errorf("field with key '%s' not found", fieldKey))
		return errs
	}

	if load, err := c.shouldLoadField(field, fieldSetKey); err != nil {
		errs = append(errs, err)
		return errs
	} else if !load {
		errs = append(errs, fmt.Errorf("field load-conditions not met"))
		return errs
	}

	for _, loader := range c.loaders {
		value, found := loader.Get(fieldSetKey, fieldKey)
		if !found {
			continue
		}

		if err := field.set(loader.Name(), value); err != nil {
			errs = append(errs, fmt.Errorf("field '%s' load error: %w", fieldKey, err))
		}
	}

	return nil
}

func (c *AppConfig) SetField(fieldSetKey, fieldKey string, fieldValue any) error {
	fieldSet, fieldSetFound := c.fieldSets[fieldSetKey]
	if !fieldSetFound {
		return fmt.Errorf("field-set with key '%s' not found", fieldSetKey)
	}

	field, fieldKeyFound := fieldSet.fieldMap[fieldKey]
	if !fieldKeyFound {
		return fmt.Errorf("field with key '%s' not found", fieldKey)
	}

	if err := field.setOverride(fieldValue); err != nil {
		return fmt.Errorf("problem setting field value: %w", err)
	}

	return nil
}

// Register loads all defined field sets and optionally checks for and handles the help flag -h and --help.
func (c *AppConfig) Register(handleHelpFlag bool) []error {
	if handleHelpFlag && len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		c.printHelpString()
		os.Exit(0)
	}

	errs := []error{}

	for _, fieldSet := range c.orderedFieldSets {
		if fieldSetErrs := c.loadFieldSet(fieldSet.Key); len(fieldSetErrs) > 0 {
			errs = append(errs, fieldSetErrs...)
			return errs
		}
	}

	c.registered = true

	return nil
}

func (c *AppConfig) HelpString() string {
	builder := strings.Builder{}

	if c.appName != "" {
		builder.WriteString(fmt.Sprintf("Usage of '%s':\n", c.appName))
	} else {
		builder.WriteString(fmt.Sprintf("Usage of '%s':\n", os.Args[0]))
	}

	if c.appDescription != "" {
		builder.WriteString(fmt.Sprintf("%s\n\n", c.appDescription))
	}

	c.addFieldsToBuilder(&builder)

	return builder.String()
}

func (c *AppConfig) ConfigMap() map[string]map[string]any {
	configMap := map[string]map[string]any{}

	for _, fieldSet := range c.fieldSets {
		fieldSetMap := map[string]any{}

		for _, field := range fieldSet.fieldMap {
			val, err := field.getValue()

			if err != nil {
				continue
			}

			if field.Sensitive {
				fieldSetMap[field.Key] = "<sensitive-value>"
				continue
			}

			if field.Type == Duration {
				val = val.(time.Duration).Milliseconds()
				fieldSetMap[fmt.Sprintf("%s_ms", field.Key)] = val

				continue
			}

			fieldSetMap[field.Key] = val
		}

		configMap[fieldSet.Key] = fieldSetMap
	}

	return configMap
}

func (c *AppConfig) GetFieldSetKeys() []string {
	keys := make([]string, len(c.fieldSets))
	idx := 0

	for key := range c.fieldSets {
		keys[idx] = key
		idx++
	}

	return keys
}

func (c *AppConfig) GetFieldSetFieldKeys(fieldSetKey string) ([]string, error) {
	fieldSet, found := c.fieldSets[fieldSetKey]
	if !found {
		return nil, fmt.Errorf("field-set not found with key: '%s'", fieldSetKey)
	}

	keys := make([]string, len(fieldSet.fieldMap))
	idx := 0

	for key := range c.fieldSets {
		keys[idx] = key
		idx++
	}

	return keys, nil
}

func (c *AppConfig) GetField(fieldSetKey, fieldKey string) (*Field, error) {
	fieldSet, found := c.fieldSets[fieldSetKey]
	if !found {
		return nil, fmt.Errorf("field-set not found with key '%s'", fieldSetKey)
	}

	field, found := fieldSet.fieldMap[fieldKey]
	if !found {
		return nil, fmt.Errorf("field-set field not found with key '%s'", fieldKey)
	}

	return field, nil
}

func (c *AppConfig) GetString(fieldSetKey, fieldKey string) (string, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, String)
	if err != nil {
		return "", err
	}

	val, _ := fieldValue.(string)

	return val, nil
}

func (c *AppConfig) GetStrings(fieldSetKey, fieldKey string) ([]string, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Strings)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]string)

	return val, nil
}

func (c *AppConfig) GetInt(fieldSetKey, fieldKey string) (int, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Int)
	if err != nil {
		return 0, err
	}

	val, _ := fieldValue.(int)

	return val, nil
}

func (c *AppConfig) GetInts(fieldSetKey, fieldKey string) ([]int, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Ints)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]int)

	return val, nil
}

func (c *AppConfig) GetBool(fieldSetKey, fieldKey string) (bool, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Bool)
	if err != nil {
		return false, err
	}

	val, _ := fieldValue.(bool)

	return val, nil
}

func (c *AppConfig) GetBools(fieldSetKey, fieldKey string) ([]bool, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Bools)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]bool)

	return val, nil
}

func (c *AppConfig) GetTime(fieldSetKey, fieldKey string) (time.Time, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Time)
	if err != nil {
		return time.Time{}, err
	}

	val, _ := fieldValue.(time.Time)

	return val, nil
}

func (c *AppConfig) GetTimes(fieldSetKey, fieldKey string) ([]time.Time, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Times)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]time.Time)

	return val, nil
}

func (c *AppConfig) GetDuration(fieldSetKey, fieldKey string) (time.Duration, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Duration)
	if err != nil {
		return 0, err
	}

	val, _ := fieldValue.(time.Duration)

	return val, nil
}

func (c *AppConfig) GetDurations(fieldSetKey, fieldKey string) ([]time.Duration, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Durations)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]time.Duration)

	return val, nil
}

func (c *AppConfig) FillStruct(configStruct any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("problem filling struct: %s", r)
		}
	}()

	if reflect.TypeOf(configStruct).Kind() != reflect.Pointer {
		return fmt.Errorf("FillStruct expects a pointer to a struct, found '%s'", reflect.TypeOf(configStruct).Kind())
	}

	configStructValue := reflect.Indirect(reflect.ValueOf(configStruct))
	configStructType := configStructValue.Type()

	if configStructValue.Kind() != reflect.Struct {
		return fmt.Errorf("FillStruct expects a pointer to a struct, found pointer to '%s'", configStructValue.Kind())
	}

	configStructField := configStructValue.FieldByName("ConfigStruct")
	if !configStructField.IsValid() || configStructField.Type().PkgPath() != "github.com/rheisen/bconf" {
		return fmt.Errorf("FillStruct expects a struct with a bconf.ConfigStruct field, none found")
	}

	configStructFieldType, _ := configStructType.FieldByName("ConfigStruct")

	baseFieldSet := configStructFieldType.Tag.Get("bconf")

	if overrideValue := configStructField.FieldByName("FieldSet"); overrideValue.String() != "" {
		baseFieldSet = overrideValue.String()
	}

	for i := 0; i < configStructValue.NumField(); i++ {
		field := configStructType.Field(i)

		if field.Name == "ConfigStruct" && field.Type.PkgPath() == "github.com/rheisen/bconf" {
			continue
		}

		fieldTagValue := field.Tag.Get("bconf")
		fieldKey := ""
		fieldSetKey := baseFieldSet

		switch fieldTagValue {
		case "":
			fieldKey = field.Name
		case "-":
			continue
		default:
			fieldTagParams := strings.Split(fieldTagValue, ",")
			fieldLocation := strings.Split(fieldTagParams[0], ".")

			fieldKey = fieldLocation[0]

			// NOTE: error if fieldLocation format isn't <field>.<field-name> ?
			if len(fieldLocation) > 1 {
				fieldSetKey = fieldLocation[0]
				fieldKey = fieldLocation[1]
			}
		}

		if fieldSetKey == "" {
			return fmt.Errorf("unidentified field-set for field: %s", fieldKey)
		}

		appConfigField, err := c.GetField(fieldSetKey, fieldKey)
		if err != nil {
			return fmt.Errorf("problem getting field '%s.%s': %w", fieldSetKey, fieldKey, err)
		}

		val, err := appConfigField.getValue()
		if err != nil && err.Error() == emptyFieldError {
			continue
		} else if err != nil {
			return fmt.Errorf("problem getting field '%s.%s' value: %w", fieldSetKey, fieldKey, err)
		}

		configStructValue.Field(i).Set(reflect.ValueOf(val))
	}

	return nil
}

// -- Private methods --

func (c *AppConfig) addFieldSet(fieldSet *FieldSet, lock bool) []error {
	if lock {
		c.fieldSetLock.Lock()
		defer c.fieldSetLock.Unlock()
	}

	fieldSet = fieldSet.Clone()

	if errs := c.checkForFieldSetStructuralIntegrity(fieldSet); len(errs) > 0 {
		return errs
	}

	if _, keyFound := c.fieldSets[fieldSet.Key]; keyFound {
		return []error{fmt.Errorf("duplicate field-set key found: '%s'", fieldSet.Key)}
	}

	fieldSet.initializeFieldMap()

	if errs := c.checkForFieldSetDependencies(fieldSet); len(errs) > 0 {
		return errs
	}

	if errs := c.generateFieldSetDefaultValues(fieldSet); len(errs) > 0 {
		return errs
	}

	if errs := c.checkForFieldSetFieldsValidity(fieldSet); len(errs) > 0 {
		return errs
	}

	fieldSet.Fields = nil

	c.fieldSets[fieldSet.Key] = fieldSet
	c.orderedFieldSets = append(c.orderedFieldSets, fieldSet)

	return nil
}

func (c *AppConfig) checkForFieldSetStructuralIntegrity(fieldSet *FieldSet) []error {
	errs := []error{}

	if fieldSetErrs := fieldSet.validate(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(errs, fmt.Errorf("field-set '%s' validation error: %w", fieldSet.Key, err))
		}
	}

	return errs
}

func (c *AppConfig) checkForFieldSetDependencies(fieldSet *FieldSet) []error {
	errs := []error{}

	for _, loadCondition := range fieldSet.LoadConditions {
		fieldSetKey, fieldKey := loadCondition.FieldDependency()
		if fieldSetKey == "" && fieldKey == "" {
			continue
		}

		fieldSetDependency, found := c.fieldSets[fieldSetKey]
		if !found {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field-set dependency not found: %s", fieldSet.Key, fieldSetKey),
			)

			continue
		}

		_, found = fieldSetDependency.fieldMap[fieldKey]
		if !found {
			errs = append(
				errs,
				fmt.Errorf(
					"field-set '%s' field-set dependency field not found: %s_%s",
					fieldSet.Key, fieldSetKey, fieldKey,
				),
			)
		}
	}

	for _, field := range fieldSet.Fields {
		if err := c.checkForFieldDependencies(field, fieldSet); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (c *AppConfig) checkForFieldDependencies(field *Field, parent *FieldSet) error {
	if len(field.LoadConditions) > 0 {
		for _, loadCondition := range field.LoadConditions {
			var fieldSetDependency *FieldSet

			var found bool

			fieldSetKey, fieldKey := loadCondition.FieldDependency()

			if fieldSetKey == "" && fieldKey == "" {
				continue
			}

			if fieldSetKey == "" || fieldSetKey == parent.Key {
				fieldSetKey = parent.Key
				fieldSetDependency = parent
			} else {
				fieldSetDependency, found = c.fieldSets[fieldSetKey]

				if !found {
					return fmt.Errorf(
						"field-set '%s' field '%s' field-set dependency not found: %s",
						fieldSetKey, field.Key, fieldSetKey,
					)
				}
			}

			if _, found = fieldSetDependency.fieldMap[fieldKey]; !found {
				return fmt.Errorf(
					"field-set '%s' field '%s' field-set field not found: %s_%s",
					parent.Key, field.Key, fieldSetKey, fieldKey,
				)
			}
		}
	}

	return nil
}

func (c *AppConfig) generateFieldSetDefaultValues(fieldSet *FieldSet) []error {
	errs := []error{}

	if fieldSetErrs := fieldSet.generateFieldDefaults(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field default value generation error: %w", fieldSet.Key, err),
			)
		}
	}

	return errs
}

func (c *AppConfig) checkForFieldSetFieldsValidity(fieldSet *FieldSet) []error {
	errs := []error{}

	if fieldSetErrs := fieldSet.validateFields(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field validation error: %w", fieldSet.Key, err),
			)
		}
	}

	return errs
}

func (c *AppConfig) loadFieldSet(fieldSetKey string) []error {
	errs := []error{}

	fieldSet, fieldSetFound := c.fieldSets[fieldSetKey]
	if !fieldSetFound {
		errs = append(errs, fmt.Errorf("field-set with key '%s' not found", fieldSetKey))
		return errs
	}

	if load, err := c.shouldLoadFieldSet(fieldSet); err != nil {
		return append(errs, err)
	} else if !load {
		return errs
	}

	for _, loader := range c.loaders {
		values := loader.GetMap(fieldSetKey, c.fieldSets[fieldSetKey].fieldKeys())
		for key, value := range values {
			field := c.fieldSets[fieldSetKey].fieldMap[key]

			if load, err := c.shouldLoadField(field, fieldSetKey); err != nil {
				errs = append(errs, err)
				continue
			} else if !load {
				continue
			}

			if err := c.fieldSets[fieldSetKey].fieldMap[key].set(loader.Name(), value); err != nil {
				errs = append(errs, fmt.Errorf("field '%s' load error: %w", key, err))
			}
		}
	}

	for _, field := range fieldSet.fieldMap {
		if field.Required && len(field.LoadConditions) < 1 {
			if _, err := field.getValue(); err != nil {
				errs = append(errs, fmt.Errorf("required field '%s_%s' not set", fieldSet.Key, field.Key))
			}
		} else if field.Required {
			if load, _ := c.shouldLoadField(field, fieldSet.Key); load {
				if _, err := field.getValue(); err != nil {
					errs = append(errs, fmt.Errorf(
						"conditionally required field '%s_%s' load condition met, but field value not set",
						fieldSet.Key,
						field.Key,
					))
				}
			}
		}
	}

	return errs
}

func (c *AppConfig) shouldLoadFieldSet(fieldSet *FieldSet) (bool, error) {
	loadFieldSet := true

	if len(fieldSet.LoadConditions) > 0 {
		for _, loadCondition := range fieldSet.LoadConditions {
			if !loadFieldSet {
				break
			}

			conditionFieldSetKey, conditionFieldSetFieldKey := loadCondition.FieldDependency()
			if conditionFieldSetKey != "" && conditionFieldSetFieldKey != "" {
				fieldValue, err := c.getFieldValue(conditionFieldSetKey, conditionFieldSetFieldKey, "any")
				if err != nil {
					return false, fmt.Errorf("problem getting field value for load condition: %w", err)
				}

				loadFieldSet, err = loadCondition.Load(fieldValue)
				if err != nil {
					return false, fmt.Errorf("problem getting load condition outcome: %w", err)
				}

				continue
			}

			var err error

			loadFieldSet, err = loadCondition.Load(nil)
			if err != nil {
				return false, fmt.Errorf("problem getting load condition outcome: %w", err)
			}

			continue
		}
	}

	return loadFieldSet, nil
}

func (c *AppConfig) shouldLoadField(field *Field, fieldSetKey string) (bool, error) {
	loadField := true

	if len(field.LoadConditions) > 0 {
		for _, loadCondition := range field.LoadConditions {
			if !loadField {
				break
			}

			conditionFieldSetKey, conditionFieldSetFieldKey := loadCondition.FieldDependency()
			if conditionFieldSetKey == "" {
				conditionFieldSetKey = fieldSetKey
			}

			if conditionFieldSetKey != "" && conditionFieldSetFieldKey != "" {
				fieldValue, err := c.getFieldValue(conditionFieldSetKey, conditionFieldSetFieldKey, "any")
				if err != nil {
					return false, fmt.Errorf("problem getting field value for load condition: %w", err)
				}

				loadField, err = loadCondition.Load(fieldValue)
				if err != nil {
					return false, fmt.Errorf("problem getting load condition outcome: %w", err)
				}

				continue
			}

			var err error

			loadField, err = loadCondition.Load(nil)
			if err != nil {
				return false, fmt.Errorf("problem getting load condition outcome: %w", err)
			}

			continue
		}
	}

	return loadField, nil
}

func (c *AppConfig) getFieldValue(fieldSetKey, fieldKey, expectedType string) (any, error) {
	field, err := c.GetField(fieldSetKey, fieldKey)
	if err != nil {
		return nil, err
	}

	if expectedType != "" && expectedType != "any" && field.Type != expectedType {
		return nil, fmt.Errorf("incorrect field-type for field '%s', found '%s'", fieldKey, field.Type)
	}

	fieldValue, err := field.getValue()
	if err != nil {
		return nil, fmt.Errorf("no value set for field '%s'", fieldKey)
	}

	return fieldValue, nil
}

type fieldEntry struct {
	fieldSetKey    string
	field          *Field
	loadConditions LoadConditions
}

func (c *AppConfig) fields() map[string]*fieldEntry {
	fields := map[string]*fieldEntry{}

	for fieldSetKey, fieldSet := range c.fieldSets {
		for _, field := range fieldSet.fieldMap {
			entry := fieldEntry{field: field, fieldSetKey: fieldSetKey}

			if len(fieldSet.LoadConditions) > 0 {
				entry.loadConditions = fieldSet.LoadConditions
			}

			if len(field.LoadConditions) > 0 {
				entry.loadConditions = append(entry.loadConditions, field.LoadConditions...)
			}

			fields[fmt.Sprintf("%s_%s", fieldSetKey, field.Key)] = &entry
		}
	}

	return fields
}

func (c *AppConfig) addFieldsToBuilder(builder *strings.Builder) {
	fields := c.fields()
	if len(fields) > 0 {
		keys := make([]string, len(fields))
		idx := 0

		for key := range fields {
			keys[idx] = key
			idx++
		}

		sort.Strings(keys)

		conditionallyRequiredFields := []string{}
		requiredFields := []string{}
		optionalFields := []string{}

		for _, key := range keys {
			fieldEntry := fields[key]

			switch {
			case fieldEntry.field.Required && fieldEntry.loadConditions == nil:
				requiredFields = append(requiredFields, key)
			case fieldEntry.field.Required && fieldEntry.loadConditions != nil:
				conditionallyRequiredFields = append(conditionallyRequiredFields, key)
			default:
				optionalFields = append(optionalFields, key)
			}
		}

		if len(requiredFields) > 0 {
			builder.WriteString("Required Configuration:\n")

			for _, key := range requiredFields {
				fmt.Fprintf(builder, "\t%s", c.fieldHelpString(fields, key))
			}
		}

		if len(conditionallyRequiredFields) > 0 {
			builder.WriteString("Conditionally Required Configuration:\n")

			for _, key := range conditionallyRequiredFields {
				fmt.Fprintf(builder, "\t%s", c.fieldHelpString(fields, key))
			}
		}

		if len(optionalFields) > 0 {
			builder.WriteString("Optional Configuration:\n")

			for _, key := range optionalFields {
				fmt.Fprintf(builder, "\t%s", c.fieldHelpString(fields, key))
			}
		}
	}
}

func (c *AppConfig) fieldHelpString(fields map[string]*fieldEntry, key string) string {
	entry := fields[key]
	field := entry.field
	loadConditions := entry.loadConditions

	if field == nil {
		return "no field matching key"
	}

	builder := strings.Builder{}
	spaceBuffer := "\t\t"

	builder.WriteString(fmt.Sprintf("%s %s\n", key, field.Type))

	if field.Description != "" {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("%s\n", field.Description))
	}

	if len(field.Enumeration) > 0 {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Accepted values: %s\n", field.enumerationString()))
	}

	if field.Default != nil && field.Sensitive {
		builder.WriteString(spaceBuffer)
		builder.WriteString("Default value: '<sensitive-value>'\n")
	} else if field.Default != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Default value: '%v'\n", field.Default))
	}

	if field.DefaultGenerator != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString("Default value: <generated-at-run-time>\n")
	}

	for _, loader := range c.loaders {
		helpString := loader.HelpString(entry.fieldSetKey, entry.field.Key)
		if helpString != "" {
			builder.WriteString(spaceBuffer)
			builder.WriteString(fmt.Sprintf("%s\n", helpString))
		}
	}

	for _, condition := range loadConditions {
		fieldSetDependency, fieldDependency := condition.FieldDependency()
		if fieldSetDependency != "" && fieldDependency != "" {
			builder.WriteString(spaceBuffer)
			builder.WriteString(
				fmt.Sprintf("Loading depends on field: '%s_%s'\n", fieldSetDependency, fieldDependency),
			)
		} else {
			builder.WriteString(spaceBuffer)
			builder.WriteString("Loading depends on: <custom-load-condition-function>\n")
		}
	}

	return builder.String()
}

func (c *AppConfig) printHelpString() {
	fmt.Printf("%s", c.HelpString())
}

// func (c *AppConfig) fieldSetLoadOrder() ([]*FieldSet, error) {
// 	fieldSets := make([]*FieldSet, len(c.fieldSets))
// 	fieldSetAvailable := map[string]struct{}{}

// 	var iter func(fieldSet *FieldSet, seen map[string]struct{}) error
// 	iter = func(fieldSet *FieldSet, seen map[string]struct{}) error {
// 		if _, seen := seen[fieldSet.Key]; seen {
// 			return fmt.Errorf("field-set cycle detected")
// 		}

// 		if _, available := fieldSetAvailable[fieldSet.Key]; available {
// 			return nil
// 		}

// 		if len(fieldSet.LoadConditions) == 0 {
// 			fieldSets = append(fieldSets, fieldSet)
// 			fieldSetAvailable[fieldSet.Key] = struct{}{}

// 			return nil
// 		}

// 		for _, condition := range fieldSet.LoadConditions {
// 			fieldSetKey, _ := condition.FieldDependency()
// 			if fieldSetKey != "" {
// 				if _, fieldSetExists := c.fieldSets[fieldSetKey]; !fieldSetExists {
// 					return fmt.Errorf("field-set dependency on non-existent field-set: '%s'", fieldSetKey)
// 				}

// 				_, fieldSetAvailable := fieldSetAvailable[fieldSetKey]
// 				if !fieldSetAvailable {
// 					seen[fieldSet.Key] = struct{}{}
// 					if err := iter(c.fieldSets[fieldSetKey], seen); err != nil {
// 						return err
// 					}

// 					continue
// 				}
// 			}
// 		}

// 		fieldSets = append(fieldSets, fieldSet)
// 		fieldSetAvailable[fieldSet.Key] = struct{}{}

// 		return nil
// 	}

// 	for _, fieldSet := range c.fieldSets {
// 		if err := iter(fieldSet, map[string]struct{}{}); err != nil {
// 			return fieldSets, err
// 		}
// 	}

// 	return fieldSets, nil
// }
