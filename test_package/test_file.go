package test_package

//go:generate main
type Enum int

const (
	EnumValue Enum = iota
	EnumValue2
	EnumValue3
	EnumValue4
)

const (
	EnumValue5, EnumValue6 Enum = 5, 6
)

type Test struct {
	Value string
}
