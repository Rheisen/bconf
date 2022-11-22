package bconfconst

const (
	Bool   = "bool"
	String = "string"
	// Uint     = "uint"
	// Uint16   = "uint16"
	// Uint32   = "uint32"
	// Uint64   = "uint64"
	Int = "int"
	// Int16    = "int16"
	// Int32    = "int32"
	// Int64    = "int64"
	// Float32  = "float32"
	Float64  = "float64"
	Time     = "time.Time"
	Duration = "time.Duration"
)

func FieldTypes() []string {
	return []string{
		Bool,
		String,
		// Uint,
		// Uint16,
		// Uint32,
		// Uint64,
		Int,
		// Int16,
		// Int32,
		// Int64,
		// Float32,
		Float64,
		Time,
		Duration,
	}
}
