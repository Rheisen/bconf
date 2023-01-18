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

func (l *EnvironmentLoader) Clone() Loader {
	newLoader := *l
	return &newLoader
}

func (l *EnvironmentLoader) Name() string {
	return "bconf_environment"
}

func (l *EnvironmentLoader) Get(key string) (string, bool) {
	return os.LookupEnv(l.environmentKey(key))
}

// func (l *EnvironmentLoader) GetMap(keys []string) map[string]string {
// 	values := map[string]string{}

// 	return values
// }

func (l *EnvironmentLoader) HelpString(key string) string {
	return fmt.Sprintf("Environment key: '%s'", l.environmentKey(key))
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
