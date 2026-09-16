package enum

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		enum           *Enum
		forceUndefined bool
		wantErr        error
	}{
		{
			name:    "nil enum",
			enum:    nil,
			wantErr: ErrValidationNoValues,
		},
		{
			name:    "empty enum",
			enum:    &Enum{},
			wantErr: ErrValidationNoValues,
		},
		{
			name: "duplicate values",
			enum: &Enum{
				{Name: "Undefined", Value: 0},
				{Name: "Ready", Value: 0},
			},
			wantErr: ErrValidationDuplicateValue,
		},
		{
			name: "missing forced undefined",
			enum: &Enum{
				{Name: "Ready", Value: 1},
			},
			forceUndefined: true,
			wantErr:        ErrValidationNoUndefined,
		},
		{
			name: "valid without forced undefined",
			enum: &Enum{
				{Name: "Ready", Value: 1},
			},
		},
		{
			name: "valid with forced undefined",
			enum: &Enum{
				{Name: "Undefined", Value: 0},
				{Name: "Ready", Value: 1},
			},
			forceUndefined: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.enum.Validate(tt.forceUndefined)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestHasUndefined(t *testing.T) {
	if !(Enum{{Name: "Undefined"}}).HasUndefined() {
		t.Fatal("HasUndefined() = false, want true")
	}

	if (Enum{{Name: "StatusUndefined"}}).HasUndefined() {
		t.Fatal("HasUndefined() = true for a non-exact undefined name")
	}

	if (Enum{nil}).HasUndefined() {
		t.Fatal("HasUndefined() = true for nil enum data")
	}
}
