package test_package

//go:generate enum_codegen
type Enum int

const (
	EnumValue1 Enum = iota // Translate="Значение Enum 1"
	EnumValue2             // Translate="Значение Enum 2"
	EnumValue3             // Translate="Значение Enum 3"
	EnumValue4             // Translate="Значение Enum 4"
)

const (
	EnumValue5, EnumValue6 Enum = 5, 6
)

type Test struct {
	Value string
}
