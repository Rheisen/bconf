package bconf

import (
	"fmt"
	"os"
	"strings"
)

func NewEnvironmentLoader() *EnvironmentLoader {
	return NewEnvironmentLoaderWithKeyPrefix("")
}

func NewEnvironmentLoaderWithKeyPrefix(keyPrefix string) *EnvironmentLoader {
	return &EnvironmentLoader{KeyPrefix: keyPrefix}
}

type EnvironmentLoader struct {
	KeyPrefix string
}

func (l *EnvironmentLoader) Clone() *EnvironmentLoader {
	newLoader := *l
	return &newLoader
}

func (l *EnvironmentLoader) CloneLoader() Loader {
	return l.Clone()
}

func (l *EnvironmentLoader) Name() string {
	return "bconf_environment"
}

func (l *EnvironmentLoader) Get(fieldSetKey, fieldKey string) (string, bool) {
	return os.LookupEnv(l.environmentKey(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)))
}

func (l *EnvironmentLoader) GetMap(fieldSetKey string, fieldKeys []string) map[string]string {
	values := map[string]string{}

	for _, fieldKey := range fieldKeys {
		value, found := os.LookupEnv(l.environmentKey(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)))
		if found {
			values[fieldKey] = value
		}
	}

	return values
}

func (l *EnvironmentLoader) HelpString(fieldSetKey, fieldKey string) string {
	return fmt.Sprintf("Environment key: '%s'", l.environmentKey(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)))
}

func (l *EnvironmentLoader) environmentKey(key string) string {
	envKey := ""
	if l.KeyPrefix != "" {
		envKey = fmt.Sprintf("%s_%s", l.KeyPrefix, key)
	} else {
		envKey = key
	}

	return strings.ToUpper(envKey)
}
