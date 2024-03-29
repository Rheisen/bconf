package bconf

import (
	"fmt"
	"os"
	"strings"
)

func NewFlagLoader() *FlagLoader {
	return NewFlagLoaderWithKeyPrefix("")
}

func NewFlagLoaderWithKeyPrefix(keyPrefix string) *FlagLoader {
	return &FlagLoader{KeyPrefix: keyPrefix}
}

type FlagLoader struct {
	KeyPrefix      string
	OverrideLookup []string
}

func (l *FlagLoader) Clone() *FlagLoader {
	clone := *l

	if len(l.OverrideLookup) > 0 {
		_ = copy(clone.OverrideLookup, l.OverrideLookup)
	}

	return &clone
}

func (l *FlagLoader) CloneLoader() Loader {
	return l.Clone()
}

func (l *FlagLoader) Name() string {
	return "bconf_flags"
}

func (l *FlagLoader) Get(fieldSetKey, fieldKey string) (string, bool) {
	values := l.flagValues()

	value, found := values[fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)]
	if found {
		return value, true
	}

	return "", false
}

func (l *FlagLoader) GetMap(fieldSetKey string, fieldKeys []string) map[string]string {
	values := map[string]string{}

	flagValues := l.flagValues()

	for _, fieldKey := range fieldKeys {
		value, found := flagValues[fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)]
		if found {
			values[fieldKey] = value
		}
	}

	return values
}

func (l *FlagLoader) HelpString(fieldSetKey, fieldKey string) string {
	return fmt.Sprintf("Flag argument: '--%s'", l.flagKey(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)))
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

func (l *FlagLoader) flagValues() map[string]string {
	values := map[string]string{}

	var args []string

	if len(l.OverrideLookup) > 0 {
		args = l.OverrideLookup
	} else {
		args = os.Args[1:]
	}

	if len(args) < 1 {
		return values
	}

	argIdx := 0
	parseArgs := true

	for parseArgs && argIdx < len(args) {
		arg := args[argIdx]

		switch {
		case strings.HasPrefix(arg, "--"):
			arg = arg[2:]
		case strings.HasPrefix(arg, "-"):
			arg = arg[1:]
		default:
			argIdx++
			continue
		}

		flagKey := ""
		flagValue := ""

		if splitIndex := strings.Index(arg, "="); splitIndex > -1 {
			flagKey = arg[:splitIndex]
			flagValue = arg[splitIndex+1:]
			values[flagKey] = flagValue
			argIdx++

			continue
		}

		flagKey = arg

		if argIdx+1 < len(args) {
			nextArg := args[argIdx+1]

			if !strings.HasPrefix(nextArg, "--") && !strings.HasPrefix(nextArg, "-") {
				values[flagKey] = nextArg
				argIdx += 2

				continue
			}
		}

		values[flagKey] = "true"
		argIdx++

		continue
	}

	return values
}
