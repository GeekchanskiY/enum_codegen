package enum

import (
	"errors"
)

const undefined = "Undefined"

var (
	ErrValidationNoValues       = errors.New("no values found")
	ErrValidationNoUndefined    = errors.New("no undefined value found")
	ErrValidationDuplicateValue = errors.New("duplicate value found")
)

func (e *Enum) Validate(forceUndefined bool) error {
	if e == nil {
		return ErrValidationNoValues
	}

	if !e.checkNotEmpty() {
		return ErrValidationNoValues
	}

	if !e.checkNoDuplicates() {
		return ErrValidationDuplicateValue
	}

	if forceUndefined {
		if !e.checkUndefinedExists() {
			return ErrValidationNoUndefined
		}
	}

	return nil
}

func (e *Enum) checkUndefinedExists() bool {
	for _, v := range *e {
		if v.Name == undefined {
			return true
		}
	}

	return false
}

func (e *Enum) checkNoDuplicates() bool {
	for i, v1 := range *e {
		for j, v2 := range *e {
			if i != j {
				if v1.Value == v2.Value {
					return false
				}
			}
		}
	}

	return true
}

func (e *Enum) checkNotEmpty() bool {
	return len(*e) != 0
}
