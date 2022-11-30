package bconf

// NewAppConfig instantiates and initializes an app configuration as defined by the given bconf.AppConfigDefinition.
func NewAppConfig(d *AppConfigDefinition) (*AppConfig, []error) {
	configDefinition := d.Clone()
	config := &AppConfig{config: &configDefinition}

	if errs := config.initialize(); len(errs) > 0 {
		return config, errs
	}

	return config, nil
}
