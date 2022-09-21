package bconfenum

type ConfigLoader int

const (
	ConfigLoaderUndefined ConfigLoader = iota
	ConfigLoaderInvalid
	ConfigLoaderError
	ConfigLoaderEnvironment
	ConfigLoaderFlags
)

var configLoaderStringMap = map[ConfigLoader]string{
	ConfigLoaderUndefined:   "undefined_config_loader",
	ConfigLoaderInvalid:     "invalid_config_loader",
	ConfigLoaderEnvironment: "environment",
	ConfigLoaderFlags:       "flags",
}

func (e ConfigLoader) IsValid() bool {
	_, found := configLoaderStringMap[e]
	return e != ConfigLoaderUndefined && e != ConfigLoaderInvalid && found
}

func (e ConfigLoader) String() string {
	configLoaderName, found := configLoaderStringMap[e]
	if !found {
		return configLoaderStringMap[ConfigLoaderInvalid]
	}

	return configLoaderName
}
