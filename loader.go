package bconf

type Loader interface {
	Clone() Loader
	Name() string
	Get(key string) (value string, found bool)
	HelpString(key string) string
}

type LoaderKeyOverride struct {
	LoaderName     string
	KeyOverride    string
	IgnorePrefixes bool
}
