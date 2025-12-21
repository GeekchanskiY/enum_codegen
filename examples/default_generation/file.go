package default_generation

//go:generate enum_codegen
type Enum int

const (
	// Undefined Value="Undefined" Translate="Enum Value is undefined"
	Undefined Enum = iota
	// EnumValue1 docs Value="Sample value"
	EnumValue1
	EnumValue2
	// EnumValue3 Some documentation, some explanation, etc. Translate="Enum ultra value"
	EnumValue3
	// Translate="Enum last value" and some documentation here
	EnumValue4
	EnumValue5
)
