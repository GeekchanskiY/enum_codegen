package force_undefined

// Enum Throws error "no undefined value found"
//
//go:generate enum_codegen -f
type Enum int

const (
	// EnumValue1 docs Value="Sample value"
	EnumValue1 Enum = iota
	EnumValue2
	// EnumValue3 Some documentation, some explanation, etc. Translate="Enum ultra value"
	EnumValue3
	// Translate="Enum last value" and some documentation here
	EnumValue4
	EnumValue5
)
