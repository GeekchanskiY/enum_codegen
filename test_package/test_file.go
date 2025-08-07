package test_package

//go:generate enum_codegen
type Enum int

const (
	// Undefined Translate="Значение Enum не определено"
	Undefined Enum = iota
	// EnumValue1 Translate="Значение Enum 1"
	EnumValue1
	// EnumValue2 Translate="Значение Enum 2"
	EnumValue2
	// EnumValue3 Translate="Значение Enum 3"
	EnumValue3
	// Translate="Значение Enum 4"
	EnumValue4
)
