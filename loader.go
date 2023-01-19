package bconf

type Loader interface {
	CloneLoader() Loader
	Name() string
	Get(fieldSetKey, fieldKey string) (value string, found bool)
	GetMap(fieldSetKey string, fieldKeys []string) (fieldValues map[string]string)
	HelpString(fieldSetKey, fieldKey string) string
}

type LoaderKeyOverride struct {
	LoaderName     string
	KeyOverride    string
	IgnorePrefixes bool
}
