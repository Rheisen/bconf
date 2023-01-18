# `bconf`: better / builder configuration for go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/rheisen/bconf?status.svg)](https://pkg.go.dev/github.com/rheisen/bconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/rheisen/bconf)](https://goreportcard.com/report/github.com/rheisen/bconf)
[![Build Status](https://github.com/rheisen/bconf/actions/workflows/golang-test.yml/badge.svg?branch=main)](https://github.com/rheisen/bconf/actions/workflows/golang-test.yml)
[![codecov.io](https://codecov.io/github/rheisen/bconf/coverage.svg?branch=main)](https://codecov.io/github/rheisen/bconf?branch=main)

`bconf` is an opinionated configuration framework that makes it easy to define, load, and validate configuration values.

```sh
go get github.com/rheisen/bconf
```

### Example

Below is an example of a `bconf.AppConfig` defined first with builders, and then with structs. Below this code block 
the behavior of the example is discussed.

```go
configuration := bconf.NewAppConfig(
    "external_http_api",
    "HTTP API for user authentication and authorization",
)

_ := configuration.SetLoaders(
    &bconf.EnvironmentLoader{KeyPrefix: "ext_http_api"},
    &bconf.FlagLoader{},
)

_ := configruation.AddFieldSets(
    bconf.FSB().Key("app").Fields(
        bconf.FB().
            Key("id").Type(bconf.String).Description("Application identifier").
            DefaultGenerator(
                func() (any, error) {
                    return fmt.Sprintf("%s", uuid.NewV4().String()), nil
                }
            ).Create(),
        bconf.FB().
            Key("session_secret").Type(bconf.String).Description("Application secret for session management").
            Sensitive().Required().
            Validator(
                func(fieldValue any) error {
                    secret, _ := fieldValue.(string)

                    minLength := 20
                    if len(secret) < minLength {
                        return fmt.Errorf(
                            "problem setting session_secret: expected string of minimum %d characters (len=%d).",
                            minLength,
                            len(secret),
                        )
                    }
                }
            ).Create(),
    ).Create(),
    bconf.FSB().Key("log").Fields(
        bconf.FB().
            Key("level").Type(bconf.String).Description("Logging level").Default("info").
            Enumeration("debug", "info", "warn", "error").Create(),
        bconf.FB().
            Key("format").Type(bconf.String).Description("Logging format").Default("json").
            Enumeration("console", "json").Create(),
        bconf.FB().
            Key("color_enabled").Type(bconf.Bool).Description("Colored logs when format is 'console'").Default(true).
            Create(),
    ).Create(),
)

// Register with the option to handle --help / -h flag set to true
if errs := configuration.Register(true); len(errs) > 0 {
    // handle configuration load errors
}

// returns the log level found in order of: default -> environment -> flag -> user override
// (based on the loaders set above).
logLevel, err := b.GetString("log", "level") 
if err != nil {
    // handle retrieval error
}
```

```go
configuration := bconf.NewAppConfig(
    "external_http_api",
    "HTTP API for user authentication and authorization",
)

_ := configuration.SetLoaders(
    &bconf.EnvironmentLoader{KeyPrefix: "ext_http_api"},
    &bconf.FlagLoader{},
)

_ := configruation.AddFieldSets(
    &bconf.FieldSet{
        Key: "app",
        Fields: bconf.Fields{
            {
                Key: "id", 
                Type: bconf.String,
                Description: "Application identifier",
                DefaultGenerator: func () (any, error) {
                    return fmt.Sprintf("%s", uuid.NewV4().String()), nil
                },
            },
            {
                Key: "session_secret",
                Type: bconf.String,
                Description: "Application secret for session management",
                Sensitive: true,
                Required: true,
                Validator: func(fieldValue any) error {
                    secret, _ := fieldValue.(string)

                    minLength := 20
                    if len(secret) < minLength {
                        return fmt.Errorf(
                            "problem setting session_secret: expected string of minimum %d characters (len=%d).",
                            minLength,
                            len(secret),
                        )
                    }
                },
            },
        },
    },
    &bconf.FieldSet{
        Key: "log",
        Fields: bconf.Fields{
            {
                Key: "level",
                Type: bconf.String,
                Description: "Logging level",
                Default: "info",
                Enumeration: []any{"debug","info","warn","error"},
            },
            {
                Key: "format",
                Type: bconf.String,
                Description: "Logging format",
                Default: "json",
                Enumeration: []any{"console", "json"},
            },
            {
                Key: "color_enabled",
                Type: bconf.Bool,
                Description: "Colored logs when format is 'console'",
                Default: true,
            },
        }
    },
)

// Register with the option to handle --help / -h flag set to true
if errs := configuration.Register(true); len(errs) > 0 {
    // handle configuration load errors
}

// returns the log level found in order of: default -> environment -> flag -> user override
// (based on the loaders set above).
logLevel, err := b.GetString("log", "level") 
if err != nil {
    // handle retrieval error
}
```
