package bconfconst

const (
	EnvironmentLoader = "environment"
	FlagLoader        = "flags"
)

func Loaders() []string {
	return []string{
		EnvironmentLoader,
		FlagLoader,
	}
}

func LoadersMap() map[string]struct{} {
	return map[string]struct{}{
		EnvironmentLoader: {},
		FlagLoader:        {},
	}
}
