package bconf

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rheisen/bconf/bconfconst"
)

func NewAppConfig(appName string, appDescription string) *AppConfig {
	return &AppConfig{
		appName:          appName,
		appDescription:   appDescription,
		fieldSets:        map[string]*FieldSet{},
		orderedFieldSets: []*FieldSet{},
		loaders:          []Loader{},
	}
}

type AppConfig struct {
	appName          string
	appDescription   string
	fieldSets        map[string]*FieldSet
	orderedFieldSets []*FieldSet
	fieldSetLock     sync.Mutex
	loaders          []Loader
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
		clonedLoaders[index] = loader.Clone()
	}

	loaderNames := make(map[string]struct{}, len(clonedLoaders))
	for _, loader := range clonedLoaders {
		if _, found := loaderNames[loader.Name()]; found {
			errs = append(errs, fmt.Errorf("duplicate loader name found: '%s'", loader.Name()))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	c.loaders = clonedLoaders

	return nil
}

func (c *AppConfig) AddFieldSet(fieldSet *FieldSet) []error {
	c.fieldSetLock.Lock()
	defer c.fieldSetLock.Unlock()

	errs := []error{}
	fieldSet = fieldSet.Clone()

	// check for field set structural integrity
	if fieldSetErrs := fieldSet.validate(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(errs, fmt.Errorf("field-set '%s' validation error: %w", fieldSet.Key, err))
		}
		return errs
	}

	// check that field set load conditions are all present
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

	if len(errs) > 0 {
		return errs
	}

	fieldSet.initializeFieldMap()

	// generate field-set field default values
	if fieldSetErrs := fieldSet.generateFieldDefaults(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field default value generation error: %w", fieldSet.Key, err),
			)
		}
		return errs
	}

	// validate field-set fields
	if fieldSetErrs := fieldSet.validateFields(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field validation error: %w", fieldSet.Key, err),
			)
		}
		return errs
	}

	// persist the field-set to AppConfig
	if c.fieldSets == nil {
		c.fieldSets = map[string]*FieldSet{fieldSet.Key: fieldSet}
		c.orderedFieldSets = append(c.orderedFieldSets, fieldSet)
		return nil
	}

	if _, keyFound := c.fieldSets[fieldSet.Key]; keyFound {
		errs = append(
			errs,
			fmt.Errorf("duplicate field-set key found: '%s'", fieldSet.Key),
		)
		return errs
	}

	c.fieldSets[fieldSet.Key] = fieldSet
	c.orderedFieldSets = append(c.orderedFieldSets, fieldSet)

	return nil
}

func (c *AppConfig) AddField(fieldSetKey string, field *Field) []error {
	c.fieldSetLock.Lock()
	defer c.fieldSetLock.Unlock()

	errs := []error{}
	errs = append(errs, fmt.Errorf("AddField not implemented"))

	return errs
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

	for _, loader := range c.loaders {
		value, found := loader.Get(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey))
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

	fields := c.fields()
	if len(fields) > 0 {
		keys := make([]string, len(fields))
		idx := 0
		for key := range fields {
			keys[idx] = key
			idx += 1
		}
		sort.Strings(keys)

		conditionallyRequiredFields := []string{}
		requiredFields := []string{}
		optionalFields := []string{}
		for _, key := range keys {
			fieldEntry := fields[key]
			if fieldEntry.field.Required && fieldEntry.loadConditions == nil {
				requiredFields = append(requiredFields, key)
			} else if fieldEntry.field.Required && fieldEntry.loadConditions != nil {
				conditionallyRequiredFields = append(conditionallyRequiredFields, key)
			} else {
				optionalFields = append(optionalFields, key)
			}
		}

		if len(requiredFields) > 0 {
			builder.WriteString("Required Configuration:\n")
			for _, key := range requiredFields {
				builder.WriteString(fmt.Sprintf("\t%s", c.fieldHelpString(fields, key)))
			}
		}

		if len(conditionallyRequiredFields) > 0 {
			builder.WriteString("Conditionally Required Configuration:\n")
			for _, key := range conditionallyRequiredFields {
				builder.WriteString(fmt.Sprintf("\t%s", c.fieldHelpString(fields, key)))
			}
		}

		if len(optionalFields) > 0 {
			builder.WriteString("Optional Configuration:\n")
			for _, key := range optionalFields {
				builder.WriteString(fmt.Sprintf("\t%s", c.fieldHelpString(fields, key)))
			}
		}
	}

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

			if field.FieldType == bconfconst.Duration {
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
		idx += 1
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
		idx += 1
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
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.String)
	if err != nil {
		return "", err
	}

	val, ok := fieldValue.(string)
	if !ok {
		return "", fmt.Errorf("problem parsing value for key '%s'", fieldSetKey)
	}

	return val, nil
}

func (c *AppConfig) GetStrings(fieldSetKey, fieldKey string) ([]string, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Strings)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]string)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetInt(fieldSetKey, fieldKey string) (int, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Int)
	if err != nil {
		return 0, err
	}

	val, ok := fieldValue.(int)
	if !ok {
		return 0, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetInts(fieldSetKey, fieldKey string) ([]int, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Ints)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]int)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetBool(fieldSetKey, fieldKey string) (bool, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Bool)
	if err != nil {
		return false, err
	}

	val, ok := fieldValue.(bool)
	if !ok {
		return false, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetBools(fieldSetKey, fieldKey string) ([]bool, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Bools)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]bool)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetTime(fieldSetKey, fieldKey string) (time.Time, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Time)
	if err != nil {
		return time.Time{}, err
	}

	val, ok := fieldValue.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetTimes(fieldSetKey, fieldKey string) ([]time.Time, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Times)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]time.Time)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetDuration(fieldSetKey, fieldKey string) (time.Duration, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Duration)
	if err != nil {
		return 0, err
	}

	val, ok := fieldValue.(time.Duration)
	if !ok {
		return 0, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

func (c *AppConfig) GetDurations(fieldSetKey, fieldKey string) ([]time.Duration, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, bconfconst.Durations)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]time.Duration)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", fieldKey)
	}

	return val, nil
}

// -- Private methods --

func (c *AppConfig) loadFieldSet(fieldSetKey string) []error {
	errs := []error{}

	fieldSet, fieldSetFound := c.fieldSets[fieldSetKey]
	if !fieldSetFound {
		errs = append(errs, fmt.Errorf("field-set with key '%s' not found", fieldSetKey))
		return errs
	}

	// Check field set load conditions
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
					errs = append(errs, fmt.Errorf("problem getting field value for load condition: %w", err))
					return errs
				}
				loadFieldSet = loadCondition.Load(fieldValue)
				continue
			}

			loadFieldSet = loadCondition.Load(nil)
			continue
		}
	}
	if !loadFieldSet {
		return errs
	}

	for _, loader := range c.loaders {
		for key, field := range c.fieldSets[fieldSetKey].fieldMap {
			value, found := loader.Get(fmt.Sprintf("%s_%s", fieldSetKey, key))
			if found {
				if err := field.set(loader.Name(), value); err != nil {
					errs = append(errs, fmt.Errorf("field '%s' load error: %w", key, err))
				}
			}
		}
	}

	for _, field := range fieldSet.fieldMap {
		if field.Required {
			if _, err := field.getValue(); err != nil {
				errs = append(errs, fmt.Errorf("required field '%s_%s' not set", fieldSet.Key, field.Key))
			}
		}
	}

	return errs
}

func (c *AppConfig) getFieldValue(fieldSetKey, fieldKey string, expectedType string) (any, error) {
	field, err := c.GetField(fieldSetKey, fieldKey)
	if err != nil {
		return nil, err
	}

	if expectedType != "" && expectedType != "any" && field.FieldType != expectedType {
		return nil, fmt.Errorf("incorrect field-type for field '%s', found '%s'", fieldKey, field.FieldType)
	}

	fieldValue, err := field.getValue()
	if err != nil {
		return nil, fmt.Errorf("no value set for field '%s'", fieldKey)
	}

	return fieldValue, nil
}

type fieldEntry struct {
	field          *Field
	loadConditions *[]LoadCondition
}

func (c *AppConfig) fields() map[string]*fieldEntry {
	fields := map[string]*fieldEntry{}

	for fieldSetKey, fieldSet := range c.fieldSets {
		for _, field := range fieldSet.fieldMap {
			entry := fieldEntry{field: field}

			if len(fieldSet.LoadConditions) > 0 {
				entry.loadConditions = &fieldSet.LoadConditions
			}

			fields[fmt.Sprintf("%s_%s", fieldSetKey, field.Key)] = &entry
		}
	}

	return fields
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

	builder.WriteString(fmt.Sprintf("%s %s\n", key, field.FieldType))

	if field.Description != "" {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("%s\n", field.Description))
	}

	if len(field.Enumeration) > 0 {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Accepted values: %s\n", field.enumerationString()))
	}

	if field.Default != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Default value: '%v'\n", field.Default))
	}

	if field.DefaultGenerator != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString("Default value: <generated-at-run-time>\n")
	}

	for _, loader := range c.loaders {
		helpString := loader.HelpString(key)
		if helpString != "" {
			builder.WriteString(spaceBuffer)
			builder.WriteString(fmt.Sprintf("%s\n", helpString))
		}
	}

	if loadConditions != nil {
		for _, condition := range *loadConditions {
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
	}

	return builder.String()
}

func (c *AppConfig) printHelpString() {
	fmt.Printf("%s", c.HelpString())
}

func (c *AppConfig) fieldSetLoadOrder() ([]*FieldSet, error) {
	fieldSets := make([]*FieldSet, len(c.fieldSets))
	fieldSetAvailable := map[string]struct{}{}

	var iter func(fieldSet *FieldSet, seen map[string]struct{}) error
	iter = func(fieldSet *FieldSet, seen map[string]struct{}) error {
		if _, seen := seen[fieldSet.Key]; seen {
			return fmt.Errorf("field-set cycle detected")
		}

		if _, available := fieldSetAvailable[fieldSet.Key]; available {
			return nil
		}

		if len(fieldSet.LoadConditions) == 0 {
			fieldSets = append(fieldSets, fieldSet)
			fieldSetAvailable[fieldSet.Key] = struct{}{}
			return nil
		}

		for _, condition := range fieldSet.LoadConditions {
			fieldSetKey, _ := condition.FieldDependency()
			if fieldSetKey != "" {
				if _, fieldSetExists := c.fieldSets[fieldSetKey]; !fieldSetExists {
					return fmt.Errorf("field-set dependency on non-existent field-set: '%s'", fieldSetKey)
				}
				_, fieldSetAvailable := fieldSetAvailable[fieldSetKey]
				if !fieldSetAvailable {
					seen[fieldSet.Key] = struct{}{}
					if err := iter(c.fieldSets[fieldSetKey], seen); err != nil {
						return err
					}

					continue
				}
			}
		}

		fieldSets = append(fieldSets, fieldSet)
		fieldSetAvailable[fieldSet.Key] = struct{}{}

		return nil
	}

	for _, fieldSet := range c.fieldSets {
		if err := iter(fieldSet, map[string]struct{}{}); err != nil {
			return fieldSets, err
		}
	}

	return fieldSets, nil
}
