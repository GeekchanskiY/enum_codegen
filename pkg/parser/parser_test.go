package parser

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParserParseConstantsAcrossPackage(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}
	writeTestFile(t, dir, "README.md", "# ignored")
	writeTestFile(t, dir, "enum_test.go", `package sample

const TestOnly Status = 100
`)
	writeTestFile(t, dir, "base.go", `package sample

type Base int

const ExtraStatus Status = 20
`)
	target := writeTestFile(t, dir, "enum.go", `package sample

import "fmt"

//go:generate enum_codegen
type Status Base

const Other = 99

const (
	// StatusUnknown Value="unknown-status" Translate="Unknown status"
	StatusUnknown Status = iota
	StatusReady
	StatusPaused, StatusStopped Status = 10, 11
)

var DefaultStatus Status = StatusReady

func FormatStatus(status Status) string {
	return fmt.Sprint(status)
}
`)

	p, err := New(dir, target, 5)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	enumName, err := p.GetEnumName()
	if err != nil {
		t.Fatalf("GetEnumName() error = %v", err)
	}
	if enumName != "Status" {
		t.Fatalf("GetEnumName() = %q, want %q", enumName, "Status")
	}

	data, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	got := map[string]struct {
		value       int64
		snakeName   string
		translation string
	}{}
	for _, item := range data {
		got[item.Name] = struct {
			value       int64
			snakeName   string
			translation string
		}{item.Value, item.SnakeName, item.Translate}
	}

	want := map[string]struct {
		value       int64
		snakeName   string
		translation string
	}{
		"ExtraStatus":   {20, "extra_status", "extra_status"},
		"StatusUnknown": {0, "unknown-status", "Unknown status"},
		"StatusReady":   {1, "status_ready", "status_ready"},
		"StatusPaused":  {10, "status_paused", "status_paused"},
		"StatusStopped": {11, "status_stopped", "status_stopped"},
	}

	if len(got) != len(want) {
		t.Fatalf("Parse() returned %d values, want %d: %#v", len(got), len(want), got)
	}
	for name, wantValue := range want {
		if gotValue, ok := got[name]; !ok || gotValue != wantValue {
			t.Fatalf("Parse()[%s] = %#v, want %#v", name, gotValue, wantValue)
		}
	}
}

func TestParserGetEnumNameTargetNotFound(t *testing.T) {
	dir := t.TempDir()
	target := writeTestFile(t, dir, "enum.go", `package sample

//go:generate enum_codegen
type Status int

const StatusReady Status = iota
`)

	p, err := New(dir, target, 1)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = p.GetEnumName()
	if !errors.Is(err, ErrTargetNotFound) {
		t.Fatalf("GetEnumName() error = %v, want %v", err, ErrTargetNotFound)
	}

	_, err = p.Parse()
	if !errors.Is(err, ErrTargetNotFound) {
		t.Fatalf("Parse() error = %v, want %v", err, ErrTargetNotFound)
	}
}

func TestParserRejectsNonIntegerEnums(t *testing.T) {
	dir := t.TempDir()
	target := writeTestFile(t, dir, "enum.go", `package sample

//go:generate enum_codegen
type Status string

const StatusReady Status = "ready"
`)

	p, err := New(dir, target, 3)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = p.Parse()
	if !errors.Is(err, ErrParsingFailed) {
		t.Fatalf("Parse() error = %v, want %v", err, ErrParsingFailed)
	}
}

func TestNewErrors(t *testing.T) {
	t.Run("missing directory", func(t *testing.T) {
		_, err := New(filepath.Join(t.TempDir(), "missing"), "enum.go", 3)
		if err == nil {
			t.Fatal("New() error = nil, want error")
		}
	})

	t.Run("invalid source", func(t *testing.T) {
		dir := t.TempDir()
		target := writeTestFile(t, dir, "enum.go", `package sample

func broken(
`)
		_, err := New(dir, target, 3)
		if err == nil {
			t.Fatal("New() error = nil, want error")
		}
	})

	t.Run("type check error", func(t *testing.T) {
		dir := t.TempDir()
		target := writeTestFile(t, dir, "enum.go", `package sample

var _ MissingType

//go:generate enum_codegen
type Status int

const StatusReady Status = iota
`)
		_, err := New(dir, target, 5)
		if err == nil {
			t.Fatal("New() error = nil, want error")
		}
	})
}

func TestParserErrors(t *testing.T) {
	t.Run("cached enum name without enum type", func(t *testing.T) {
		p := &enumParser{enumName: "Status"}

		_, err := p.Parse()
		if !errors.Is(err, ErrTargetNotFound) {
			t.Fatalf("Parse() error = %v, want %v", err, ErrTargetNotFound)
		}
	})

	t.Run("target file missing from package", func(t *testing.T) {
		dir := t.TempDir()
		writeTestFile(t, dir, "enum.go", `package sample

type Status int
`)

		_, err := New(dir, filepath.Join(dir, "missing.go"), 3)
		if !errors.Is(err, ErrTargetNotFound) {
			t.Fatalf("New() error = %v, want %v", err, ErrTargetNotFound)
		}
	})

	t.Run("no enum constants", func(t *testing.T) {
		dir := t.TempDir()
		target := writeTestFile(t, dir, "enum.go", `package sample

//go:generate enum_codegen
type Status int

const Other = 1
`)
		p, err := New(dir, target, 3)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		_, err = p.Parse()
		if !errors.Is(err, ErrTargetNotFound) {
			t.Fatalf("Parse() error = %v, want %v", err, ErrTargetNotFound)
		}
	})

	t.Run("integer value outside int64", func(t *testing.T) {
		dir := t.TempDir()
		target := writeTestFile(t, dir, "enum.go", `package sample

//go:generate enum_codegen
type Status uint64

const StatusHuge Status = 1 << 63
`)
		p, err := New(dir, target, 3)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		_, err = p.Parse()
		if !errors.Is(err, ErrParsingFailed) {
			t.Fatalf("Parse() error = %v, want %v", err, ErrParsingFailed)
		}
	})
}

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", name, err)
	}

	return path
}
