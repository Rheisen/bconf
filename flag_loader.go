package bconf

import (
	"flag"
	"fmt"
	"strings"
)

type FlagLoader struct {
	KeyPrefix string
	values    map[string]string
}

func (l *FlagLoader) Clone() Loader {
	return l
}

func (l *FlagLoader) Name() string {
	return "bconf_flags"
}

func (l *FlagLoader) Get(key string) (string, bool) {
	value := flag.Lookup(l.flagKey(key))
	if value != nil {
		return value.Value.String(), true
	}

	return "", false
}

// func (l *FlagLoader) GetMap(keys []string) map[string]string {
// 	values := map[string]string{}

// 	return values
// }

func (l *FlagLoader) HelpString(key string) string {
	return fmt.Sprintf("Flag argument: '--%s'", l.flagKey(key))
}

func (l *FlagLoader) flagKey(key string) string {
	flagKey := ""
	if l.KeyPrefix != "" {
		flagKey = fmt.Sprintf("%s_%s", l.KeyPrefix, key)
	} else {
		flagKey = key
	}

	return strings.ToLower(flagKey)
}
