package bconfv1

type Loader interface {
	Clone() Loader
	Name() string
	Get(key string) (value string, found bool)
	// GetMap(keys []string) (values map[string]string)
	HelpString(key string) string
}

type LoaderKeyOverride struct {
	LoaderName     string
	KeyOverride    string
	IgnorePrefixes bool
}
