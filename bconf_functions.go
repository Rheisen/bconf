package bconf

// NewAppConfig instantiates and initializes an app configuration as defined by the given bconf.AppConfigDefinition.
// NewAppConfig will return all errors it identifies in two stages:
// a) Invalid Config Definition -- occurs when config fields are not typed correctly or are type inconsistent.
// b) Invalid Config Fields -- occurs when config fields do not pass requirements / enumeration checks / validation.
func NewAppConfig(d AppConfigDefinition) (*AppConfig, []error) {
	configDefinition := d.clone()
	config := &AppConfig{config: &configDefinition}
	if errs := config.initialize(); len(errs) > 0 {
		return config, errs
	}

	return config, nil
}
