// Code generated via enum_codegen DO NOT EDIT.
package test_package

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	_ sql.Scanner    = (*Enum)(nil)
	_ driver.Valuer  = (*Enum)(nil)
	_ fmt.Stringer   = (*Enum)(nil)
	_ json.Marshaler = (*Enum)(nil)
)

var Tags = map[Enum]string{
	Undefined:  "undefined",
	EnumValue1: "enum_value_1",
	EnumValue2: "enum_value_2",
	EnumValue3: "enum_value_3",
	EnumValue4: "enum_value_4",
	EnumValue5: "enum_value_5",
}

var Types = map[string]Enum{
	"undefined":    Undefined,
	"enum_value_1": EnumValue1,
	"enum_value_2": EnumValue2,
	"enum_value_3": EnumValue3,
	"enum_value_4": EnumValue4,
	"enum_value_5": EnumValue5,
}

var Translations = map[Enum]string{
	Undefined:  "Enum Value is undefined",
	EnumValue1: "enum_value_1",
	EnumValue2: "enum_value_2",
	EnumValue3: "Enum ultra value",
	EnumValue4: "Enum last value",
	EnumValue5: "enum_value_5",
}

func (t *Enum) Scan(src any) error {
	value, ok := src.(string)

	if !ok {
		return errors.New("src is not string")
	}

	*t = Undefined
	if v, ok := Types[value]; ok {
		*t = v
	}

	return nil
}

func (t Enum) Value() (driver.Value, error) {
	return Tags[t], nil
}

func (t Enum) String() string {
	return Tags[t]
}

func (t Enum) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Enum) UnmarshalJSON(data []byte) error {
	var (
		s   string
		err error
	)

	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}

	if v, ok := Types[s]; ok {
		*t = v

		return nil
	}

	return fmt.Errorf("invalid status: %s", s)
}
