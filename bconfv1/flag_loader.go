package bconf

import (
	"fmt"
	"os"
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
	values := l.flagValues()

	value, found := values[key]
	if found {
		return value, true
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

func (l *FlagLoader) flagValues() map[string]string {
	values := map[string]string{}

	args := os.Args[1:]
	if len(args) < 1 {
		return values
	}

	argIdx := 0
	parseArgs := true

	for parseArgs && argIdx < len(args) {
		arg := args[argIdx]

		if strings.HasPrefix(arg, "--") {
			arg = arg[2:]
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
		} else {
			argIdx += 1
			continue
		}

		flagKey := ""
		flagValue := ""

		if splitIndex := strings.Index(arg, "="); splitIndex > -1 {
			flagKey = arg[:splitIndex]
			flagValue = arg[splitIndex+1:]

			values[flagKey] = flagValue

			argIdx += 1
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

		argIdx += 1
		continue
	}

	return values
}
