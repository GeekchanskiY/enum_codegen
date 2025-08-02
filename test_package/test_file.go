package test_package

//go:generate main enum
type Enum string

//go:generate main
type Test struct {
	Value string
}
