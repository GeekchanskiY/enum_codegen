package parser

import "testing"

func TestCommentMetadata(t *testing.T) {
	comment := `// Value="custom_value" Translate="Custom translation"`

	if got := GetValueFromComment(comment); got != "custom_value" {
		t.Fatalf("GetValueFromComment() = %q, want %q", got, "custom_value")
	}

	if got := GetTranslationFromComment(comment); got != "Custom translation" {
		t.Fatalf("GetTranslationFromComment() = %q, want %q", got, "Custom translation")
	}

	if got := GetValueFromComment("// no metadata"); got != "" {
		t.Fatalf("GetValueFromComment() = %q, want empty string", got)
	}

	if got := GetTranslationFromComment("// no metadata"); got != "" {
		t.Fatalf("GetTranslationFromComment() = %q, want empty string", got)
	}
}

func TestCamelToSnake(t *testing.T) {
	tests := map[string]string{
		"EnumValue1":     "enum_value_1",
		"HTTPStatusCode": "httpstatus_code",
		"Value2Fast":     "value_2_fast",
		"Already_Snake":  "already_snake",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := CamelToSnake(input); got != want {
				t.Fatalf("CamelToSnake(%q) = %q, want %q", input, got, want)
			}
		})
	}
}
