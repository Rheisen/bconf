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
configuration := bconf.NewAppConfig("external_http_api", "HTTP API for user authentication and authorization")
_ := configuration.SetLoaders(
    &bconf.EnvironmentLoader{KeyPrefix: "ext_http_api"},
    &bconf.FlagLoader{},
},
_ := configruation.AddFieldSet( 
    "app",
    []*bconf.Field{
        {
            Key: "id", 
            FieldType: bconfconst.String,
            Description: "Application identifier",
            DefaultGenerator: func () (any, error) {
                return fmt.Sprintf("%s", uuid.NewV4().String()), nil
            },
        },
        {
            Key: "session_secret",
            FieldType: bconfconst.String,
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
)
_ := configuration.AddFieldSet(
    "log",
    []*bconf.Field{
        "level": {
            FieldType: bconfconst.String,
            Description: "Logging level",
            Default: "info",
            Enumeration: []any{"debug","info","warn","error"},
        },
        "format": {
            FieldType: bconfconst.String,
            Description: "Logging format",
            Default: "json",
            Enumeration: []any{"console", "json"},
        },
        "color_enabled": bconf.Field{
            FieldType: bconfconst.Bool,
            Description: "Colored logs when format is 'console'",
            Default: true,
        },
    },
)

if errs := configuration.Register(true); len(errs) > 0 {
    // handle configuration load errors
}

logLevel, err := b.GetString("log", "level") // returns the log level found in order of: default -> ENV -> Flag order
if err != nil {
    // handle retrieval error
}
```
