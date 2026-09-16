// Code generated via enum_codegen DO NOT EDIT.
package default_generation

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Force interface implementation
var (
	_ sql.Scanner    = (*Enum)(nil)
	_ driver.Valuer  = (*Enum)(nil)
	_ fmt.Stringer   = (*Enum)(nil)
	_ json.Marshaler = (*Enum)(nil)
)

var EnumTags = map[Enum]string{
	Undefined:  "Undefined",
	EnumValue1: "Sample value",
	EnumValue2: "enum_value_2",
	EnumValue3: "enum_value_3",
	EnumValue4: "enum_value_4",
	EnumValue5: "enum_value_5",
}

var EnumTypes = map[string]Enum{
	"Undefined":    Undefined,
	"Sample value": EnumValue1,
	"enum_value_2": EnumValue2,
	"enum_value_3": EnumValue3,
	"enum_value_4": EnumValue4,
	"enum_value_5": EnumValue5,
}

var EnumTranslations = map[Enum]string{
	Undefined:  "Enum Value is undefined",
	EnumValue1: "enum_value_1",
	EnumValue2: "enum_value_2",
	EnumValue3: "Enum ultra value",
	EnumValue4: "Enum last value",
	EnumValue5: "enum_value_5",
}

func (t *Enum) Scan(src any) error {
	var value string
	switch v := src.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	case nil:
		*t = Undefined
		return nil
	default:
		return fmt.Errorf("src is %T, not string or []byte", src)
	}

	if v, ok := EnumTypes[value]; ok {
		*t = v
		return nil
	}
	*t = Undefined
	return nil
}

func (t Enum) Value() (driver.Value, error) {
	return EnumTags[t], nil
}

func (t Enum) String() string {
	return EnumTags[t]
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

	if v, ok := EnumTypes[s]; ok {
		*t = v

		return nil
	}

	return fmt.Errorf("invalid Enum: %s", s)
}
