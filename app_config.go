package bconf

import (
	"fmt"
	"os"
	"time"

	"github.com/rheisen/bconf/bconfconst"
)

type AppConfig struct {
	config *AppConfigDefinition
}

func (c *AppConfig) GetKeys() []string {
	keys := make([]string, len(c.config.ConfigFields))
	idx := 0
	for key := range c.config.ConfigFields {
		keys[idx] = key
		idx += 1
	}

	return keys
}

func (c *AppConfig) GetField(key string) (Field, error) {
	field, found := c.config.ConfigFields[key]
	if !found {
		return Field{}, fmt.Errorf("problem finding value for key '%s'", key)
	}

	return *field, nil
}

func (c *AppConfig) GetString(key string) (string, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.String)
	if err != nil {
		return "", err
	}

	val, ok := fieldValue.(string)
	if !ok {
		return "", fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetStrings(key string) ([]string, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Strings)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]string)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetInt(key string) (int, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Int)
	if err != nil {
		return 0, err
	}

	val, ok := fieldValue.(int)
	if !ok {
		return 0, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetInts(key string) ([]int, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Ints)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]int)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetBool(key string) (bool, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Bool)
	if err != nil {
		return false, err
	}

	val, ok := fieldValue.(bool)
	if !ok {
		return false, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetBools(key string) ([]bool, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Bools)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]bool)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetTime(key string) (time.Time, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Bool)
	if err != nil {
		return time.Time{}, err
	}

	val, ok := fieldValue.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetTimes(key string) ([]time.Time, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Times)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]time.Time)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetDuration(key string) (time.Duration, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Bool)
	if err != nil {
		return 0, err
	}

	val, ok := fieldValue.(time.Duration)
	if !ok {
		return 0, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetDurations(key string) ([]time.Duration, error) {
	fieldValue, err := c.getFieldValue(key, bconfconst.Durations)
	if err != nil {
		return nil, err
	}

	val, ok := fieldValue.([]time.Duration)
	if !ok {
		return nil, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) initialize() []error {
	if c.config.HandleHelpFlag && len(os.Args) > 1 && os.Args[1] == "--help" {
		c.config.printHelpText()
		os.Exit(0)
	}

	if errs := c.config.GenerateFieldDefaults(); len(errs) > 0 {
		return errs
	}
	if errs := c.config.Validate(); len(errs) > 0 {
		return errs
	}

	if errs := c.config.loadFields(); len(errs) > 0 {
		return errs
	}

	return nil
}

func (c *AppConfig) getFieldValue(key string, expectedType string) (any, error) {
	field, err := c.GetField(key)
	if err != nil {
		return nil, err
	}

	if field.FieldType != expectedType {
		return nil, fmt.Errorf("incorrect type for value matching key '%s'", key)
	}

	fieldValue, err := field.GetValue()
	if err != nil {
		return nil, fmt.Errorf("no set value for key '%s'", key)
	}

	return fieldValue, nil
}
