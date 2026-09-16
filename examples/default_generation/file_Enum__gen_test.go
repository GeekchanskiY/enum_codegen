package default_generation

import (
	"encoding/json"
	"testing"
)

func TestGeneratedExampleMapsAndString(t *testing.T) {
	if EnumTags[EnumValue1] != "Sample value" {
		t.Fatalf("EnumTags[EnumValue1] = %q, want %q", EnumTags[EnumValue1], "Sample value")
	}

	if EnumTypes["enum_value_2"] != EnumValue2 {
		t.Fatalf("EnumTypes[enum_value_2] = %v, want %v", EnumTypes["enum_value_2"], EnumValue2)
	}

	if EnumTranslations[EnumValue3] != "Enum ultra value" {
		t.Fatalf("EnumTranslations[EnumValue3] = %q, want %q", EnumTranslations[EnumValue3], "Enum ultra value")
	}

	if EnumValue4.String() != "enum_value_4" {
		t.Fatalf("EnumValue4.String() = %q, want %q", EnumValue4.String(), "enum_value_4")
	}
}

func TestGeneratedExampleScanValueAndJSON(t *testing.T) {
	var value Enum
	if err := value.Scan([]byte("enum_value_5")); err != nil {
		t.Fatalf("Scan([]byte) error = %v", err)
	}
	if value != EnumValue5 {
		t.Fatalf("Scan([]byte) value = %v, want %v", value, EnumValue5)
	}

	if err := value.Scan("missing"); err != nil {
		t.Fatalf("Scan(missing) error = %v", err)
	}
	if value != Undefined {
		t.Fatalf("Scan(missing) value = %v, want Undefined", value)
	}

	value = EnumValue1
	if err := value.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error = %v", err)
	}
	if value != Undefined {
		t.Fatalf("Scan(nil) value = %v, want Undefined", value)
	}

	if err := value.Scan(10); err == nil {
		t.Fatal("Scan(int) error = nil, want error")
	}

	driverValue, err := EnumValue1.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	if driverValue != "Sample value" {
		t.Fatalf("Value() = %v, want %q", driverValue, "Sample value")
	}

	data, err := json.Marshal(EnumValue2)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if string(data) != `"enum_value_2"` {
		t.Fatalf("MarshalJSON() = %s, want %q", data, `"enum_value_2"`)
	}

	var decoded Enum
	if err = json.Unmarshal([]byte(`10`), &decoded); err == nil {
		t.Fatal("UnmarshalJSON(non-string) error = nil, want error")
	}

	if err = json.Unmarshal([]byte(`"enum_value_3"`), &decoded); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if decoded != EnumValue3 {
		t.Fatalf("UnmarshalJSON() value = %v, want %v", decoded, EnumValue3)
	}

	if err = json.Unmarshal([]byte(`"missing"`), &decoded); err == nil {
		t.Fatal("UnmarshalJSON(missing) error = nil, want error")
	}
}
