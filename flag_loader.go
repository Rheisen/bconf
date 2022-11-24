package bconf

import "fmt"

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
	if l.values == nil {
		// TODO: initialize
	}

	return "", false
}

func (l *FlagLoader) HelpString(key string) string {
	return fmt.Sprintf("Flag argument: '%s'", l.flagKey(key))
}

func (l *FlagLoader) flagKey(key string) string {
	return key
}
