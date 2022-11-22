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

func (c *AppConfig) GetField(key string) (Field, error) {
	field, found := c.config.ConfigFields[key]
	if !found {
		return Field{}, fmt.Errorf("problem finding value for key '%s'", key)
	}

	return *field, nil
}

func (c *AppConfig) GetString(key string) (string, error) {
	field, err := c.GetField(key)
	if err != nil {
		return "", err
	}

	if field.FieldType != bconfconst.String {
		return "", fmt.Errorf("incorrect type for value matching key '%s'", key)
	}

	fieldValue, err := field.GetValue()
	if err != nil {
		return "", fmt.Errorf("no set value for key '%s'", key)
	}

	val, ok := fieldValue.(string)
	if !ok {
		return "", fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetInt(key string) (int, error) {
	field, err := c.GetField(key)
	if err != nil {
		return 0, err
	}

	if field.FieldType != bconfconst.Int {
		return 0, fmt.Errorf("incorrect type for value matching key '%s'", key)
	}

	fieldValue, err := field.GetValue()
	if err != nil {
		return 0, fmt.Errorf("no set value for key '%s'", key)
	}

	val, ok := fieldValue.(int)
	if !ok {
		return 0, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetBool(key string) (bool, error) {
	field, err := c.GetField(key)
	if err != nil {
		return false, err
	}

	if field.FieldType != bconfconst.Bool {
		return false, fmt.Errorf("incorrect type for value matching key '%s'", key)
	}

	fieldValue, err := field.GetValue()
	if err != nil {
		return false, fmt.Errorf("no set value for key '%s'", key)
	}

	val, ok := fieldValue.(bool)
	if !ok {
		return false, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetTime(key string) (time.Time, error) {
	field, err := c.GetField(key)
	if err != nil {
		return time.Time{}, err
	}

	if field.FieldType != bconfconst.Time {
		return time.Time{}, fmt.Errorf("incorrect type for value matching key '%s'", key)
	}

	fieldValue, err := field.GetValue()
	if err != nil {
		return time.Time{}, fmt.Errorf("no set value for key '%s'", key)
	}

	val, ok := fieldValue.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) GetDuration(key string) (time.Duration, error) {
	field, err := c.GetField(key)
	if err != nil {
		return 0, err
	}

	if field.FieldType != bconfconst.Duration {
		return 0, fmt.Errorf("incorrect type for value matching key '%s'", key)
	}

	fieldValue, err := field.GetValue()
	if err != nil {
		return 0, fmt.Errorf("no set value for key '%s'", key)
	}

	val, ok := fieldValue.(time.Duration)
	if !ok {
		return 0, fmt.Errorf("problem parsing value for key '%s'", key)
	}

	return val, nil
}

func (c *AppConfig) initialize() []error {
	if c.config.HandleHelpFlag && len(os.Args) > 0 && os.Args[1] == "--help" {
		c.config.PrintHelpText()
		os.Exit(0)
	}

	if errs := c.config.GenerateFieldDefaults(); len(errs) > 0 {
		return errs
	}
	if errs := c.config.Validate(); len(errs) > 0 {
		return errs
	}

	c.config.setDefaults()
	if errs := c.config.loadFields(); len(errs) > 0 {
		return errs
	}

	return nil
}
