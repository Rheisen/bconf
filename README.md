# `bconf`: better / builder configuration for go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/rheisen/bconf?status.svg)](https://pkg.go.dev/github.com/rheisen/bconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/rheisen/bconf)](https://goreportcard.com/report/github.com/rheisen/bconf)
[![Build Status](https://github.com/rheisen/bconf/actions/workflows/golang-test.yml/badge.svg?branch=main)](https://github.com/rheisen/bconf/actions/workflows/golang-test.yml)
[![codecov.io](https://codecov.io/github/rheisen/bconf/coverage.svg?branch=main)](https://codecov.io/github/rheisen/bconf?branch=main)

`bconf` is an opinionated configuration framework that makes it easy to define, load, and validate application 
configuration values.

```sh
go get github.com/rheisen/bconf
```

### Why `bconf`

`bconf` provides tooling to write your configuration package by package. With `bconf`, configuration lives right
alongside the code that needs it. This makes it so that configuration is more easily re-used and composible by
multiple applications, just like your packages should be.

`bconf` accomplishes this with `bconf.FieldSets`, which provide a namespace and logical grouping for related
configuration. Independent packages define their `bconf.FieldSets`, and then application executables can attach them
to a `bconf.AppConfig`, which provides a unified structure for loading and retrieving configuration values.

Check out the introductory examples below, and see if `bconf` is right for your project.

### Example

Below is an example of a `bconf.AppConfig` defined first with builders, and then with structs. Below these code blocks 
the behavior of the example is discussed.

```go
configuration := bconf.NewAppConfig(
	"external_http_api",
	"HTTP API for user authentication and authorization",
)

_ = configuration.SetLoaders(
	&bconf.EnvironmentLoader{KeyPrefix: "ext_http_api"},
	&bconf.FlagLoader{},
)

_ = configuration.AddFieldSets(
	bconf.NewFieldSetBuilder().Key("app").Fields(
		bconf.NewFieldBuilder().
			Key("id").Type(bconf.String).
			Description("Application identifier").
			DefaultGenerator(
				func() (any, error) {
					return fmt.Sprintf("%s", uuid.NewV4().String()), nil
				},
			).Create(),
		bconf.FB(). // FB() is a shorthand function for NewFieldBuilder()
				Key("session_secret").Type(bconf.String).
				Description("Application secret for session management").
				Sensitive().Required().
				Validator(
				func(fieldValue any) error {
					secret, _ := fieldValue.(string)

					minLength := 20
					if len(secret) < minLength {
						return fmt.Errorf(
							"expected string of minimum %d characters (len=%d)",
							minLength,
							len(secret),
						)
					}

					return nil
				},
			).Create(),
	).Create(),
	bconf.FSB().Key("log").Fields( // FSB() is a shorthand function for NewFieldSetBuilder()
		bconf.FB().
			Key("level").Type(bconf.String).Default("info").
			Description("Logging level").
			Enumeration("debug", "info", "warn", "error").Create(),
		bconf.FB().
			Key("format").Type(bconf.String).Default("json").
			Description("Logging format").
			Enumeration("console", "json").Create(),
		bconf.FB().
			Key("color_enabled").Type(bconf.Bool).Default(true).
			Description("Colored logs when format is 'console'").
			Create(),
	).Create(),
)

// Register with the option to handle --help / -h flag set to true
if errs := configuration.Register(true); len(errs) > 0 {
	// handle configuration load errors
}

// returns the log level found in order of: default -> environment -> flag -> user override
// (based on the loaders set above).
logLevel, err := configuration.GetString("log", "level")
if err != nil {
	// handle retrieval error
}

fmt.Printf("log-level: %s", logLevel)
```

```go
configuration := bconf.NewAppConfig(
	"external_http_api",
	"HTTP API for user authentication and authorization",
)

_ = configuration.SetLoaders(
	&bconf.EnvironmentLoader{KeyPrefix: "ext_http_api"},
	&bconf.FlagLoader{},
)

_ = configuration.AddFieldSets(
	&bconf.FieldSet{
		Key: "app",
		Fields: bconf.Fields{
			{
				Key:         "id",
				Type:        bconf.String,
				Description: "Application identifier",
				DefaultGenerator: func() (any, error) {
					return uuid.NewV4().String(), nil
				},
			},
			{
				Key:         "session_secret",
				Type:        bconf.String,
				Description: "Application secret for session management",
				Sensitive:   true,
				Required:    true,
				Validator: func(fieldValue any) error {
					secret, _ := fieldValue.(string)

					minLength := 20
					if len(secret) < minLength {
						return fmt.Errorf(
							"expected string of minimum %d characters (len=%d)",
							minLength,
							len(secret),
						)
					}

					return nil
				},
			},
		},
	},
	&bconf.FieldSet{
		Key: "log",
		Fields: bconf.Fields{
			{
				Key:         "level",
				Type:        bconf.String,
				Description: "Logging level",
				Default:     "info",
				Enumeration: []any{"debug", "info", "warn", "error"},
			},
			{
				Key:         "format",
				Type:        bconf.String,
				Description: "Logging format",
				Default:     "json",
				Enumeration: []any{"console", "json"},
			},
			{
				Key:         "color_enabled",
				Type:        bconf.Bool,
				Description: "Colored logs when format is 'console'",
				Default:     true,
			},
		},
	},
)

// Register with the option to handle --help / -h flag set to true
if errs := configuration.Register(true); len(errs) > 0 {
	// handle configuration load errors here
}

// returns the log level found in order of: default -> environment -> flag -> user override
// (based on the loaders set above).
logLevel, err := configuration.GetString("log", "level")
if err != nil {
	// handle retrieval error
}

fmt.Printf("log-level: %s", logLevel)
```

In both of the code blocks above, a `bconf.AppConfig` is defined with two field-sets (which group configuration related
to the application and logging in this case), and registered with help flag parsing.

If this code was executed in a `main()` function, it would print the log level picked up by the configuration from the
flags or run-time environment before falling back on the defined default value of "info".

If this code was executed inside the `main()` function and passed a `--help` or `-h` flag, it would print the following
output:

```
Usage of 'external_http_api':
HTTP API for user authentication and authorization

Required Configuration:
        app_session_secret string
                Application secret for session management
                Environment key: 'EXT_HTTP_API_APP_SESSION_SECRET'
                Flag argument: '--app_session_secret'
Optional Configuration:
        app_id string
                Application identifier
                Default value: <generated-at-run-time>
                Environment key: 'EXT_HTTP_API_APP_ID'
                Flag argument: '--app_id'
        log_color_enabled bool
                Colored logs when format is 'console'
                Default value: 'true'
                Environment key: 'EXT_HTTP_API_LOG_COLOR_ENABLED'
                Flag argument: '--log_color_enabled'
        log_format string
                Logging format
                Accepted values: ['console', 'json']
                Default value: 'json'
                Environment key: 'EXT_HTTP_API_LOG_FORMAT'
                Flag argument: '--log_format'
        log_level string
                Logging level
                Accepted values: ['debug', 'info', 'warn', 'error']
                Default value: 'info'
                Environment key: 'EXT_HTTP_API_LOG_LEVEL'
                Flag argument: '--log_level'
```

This is a simple example where all the configuration code is in one place, but it doesn't need to be!

To view more examples, including a real-world example showcasing how configuration can live alongside package code,
please visit [github.com/rheisen/bconf-examples](https://github.com/rheisen/bconf-examples).

### Features

* Ability to generate default configuration values with the `bconf.Field` `DefaultGenerator` function
* Flexible configuration value validation with the `bconf.Field` `Validator` function
* Load configuration values from environment and flags
* Ability to conditionally load `bconf.FieldSets` via flexible `bconf.LoadConditions`

### Limitations

* No support for loading configuration values from JSON or YAML files (although a JSON or YAML file loader can be
given to `bconf` so long as it follows the `bconf.Loader` interface)
* No support for watching / automatically updating configuration values
