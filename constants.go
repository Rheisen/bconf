package bconf

const (
	Bool      = "bool"
	Bools     = "[]bool"
	String    = "string"
	Strings   = "[]string"
	Int       = "int"
	Ints      = "[]int"
	Float     = "float64"
	Floats    = "[]float64"
	Time      = "time.Time"
	Times     = "[]time.Time"
	Duration  = "time.Duration"
	Durations = "[]time.Duration"
)

func FieldTypes() []string {
	return []string{
		Bool,
		Bools,
		String,
		Strings,
		Int,
		Ints,
		Float,
		Floats,
		Time,
		Times,
		Duration,
		Durations,
	}
}
