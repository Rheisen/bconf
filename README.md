# Bconf (Better / Base Configuration)

Bconf is an opinionated configuration framework that makes it easy to load and validate configuration values from
a variety of supported "configuration loaders", e.g. environment variables, flags, etc.

### Installing

```sh
go get github.com/rheisen/bconf
```

### Example

Below is an example of a `bconf.AppConfig`. Below this code block the behavior of the example is discussed.

```go
baseConfig, errs := bconf.NewAppConfig(
    bconf.ConfigDefinition{
        Name: "external_http_api",
        KeyPrefix: "ext_http_api",
        HandleHelpFlag: true,
        Loaders: []string{bconfconst.EnvironmentLoader, bconfconst.FlagLoader},
        ConfigFields: map[string]*bconf.Field{
            "session_secret": {
                FieldType: bconfconst.String,
                Description: "Application secret for session management",
                Required: true,
                Validator: func(fieldValue any) error {
                    secret := fieldValue.(string)

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
            "log_level": {
                FieldType: bconfconst.String,
                Description: "Application logging level",
                Default: "info",
                Enumeration: []any{"debug","info","warn","error"},
            },
            "log_format": {
                FieldType: bconfconst.String,
                Description: "Application logging format",
                Default: "json",
                Enumeration: []any{"console", "json"},
            },
            "log_color_enabled": bconf.Field{
                FieldType: bconfconst.Bool,
                Description: "Application colored logs when format is 'console'",
                Default: true,
            },
            "app_id": bconf.Field{
                FieldType: bconfconst.String,
                Description: "Application identifier",
                DefaultGenerator: func () (any, error) {
                    return fmt.Sprintf("%s", uuid.NewV4().String()), nil
                },
            },
        },
    }
)

if errs != nil && len(errs) > 0 {
    // handle configuration errors
}

appLogLevel, err := b.GetString("log_level") // returns the log level found in order of: default -> ENV -> Flag order
if err != nil {
    // handle retrieval error
}
```

In order to create a `bconf.AppConfig`, you must supply a `bconf.AppConfigDefinition`. A `bconf.AppConfigDefinition`
provides public fields that enable end-users to specify the behavior around loading application configuration values.

In the example above, the `ConfigFields` parameter of the `bconf.AppConfigDefinition` provides a convenient way to map
a configuration key to a `bonf.Field`. Similar to how a `bconf.AppConfigDefinition` allows end-users to specify the
behavior around loading configuration values, a `bconf.Field` is a structure that allows an end-user to specify the
desired behavior of a specific configuration value.

Let's break down how the `session_secret` field is defined. It has a field type of "string", the field must be loaded
from at least one of the configuration loaders, and it must pass the validator function (which checks that the
secret meets a length requirement).
