package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGeneratesCode(t *testing.T) {
	dir := t.TempDir()
	writeCLITestFile(t, dir, "enum.go", `package sample

//go:generate enum_codegen
type Status int

const (
	StatusUnknown Status = iota
	StatusReady
)
`)

	var stdout, stderr bytes.Buffer
	exitCode := run(nil, testEnv(map[string]string{
		"GOFILE":    "enum.go",
		"GOPACKAGE": "sample",
		"GOLINE":    "3",
	}), func() (string, error) {
		return dir, nil
	}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("run() exit = %d, want 0; stderr: %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "generated code to") {
		t.Fatalf("stdout = %q, want generated code message", stdout.String())
	}

	generated, err := os.ReadFile(filepath.Join(dir, "enum_Status__gen.go"))
	if err != nil {
		t.Fatalf("ReadFile(generated) error = %v", err)
	}
	if !strings.Contains(string(generated), "var StatusTags") {
		t.Fatalf("generated file does not contain StatusTags:\n%s", generated)
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"-h"}, testEnv(nil), os.Getwd, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run(-h) exit = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "Usage: enum_codegen") {
		t.Fatalf("stdout = %q, want usage", stdout.String())
	}
}

func TestRunErrors(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		getwd      func() (string, error)
		wantCode   int
		wantStderr string
	}{
		{
			name:       "invalid flag",
			args:       []string{"-unknown"},
			wantCode:   2,
			wantStderr: "flag provided but not defined",
		},
		{
			name:       "invalid goline",
			env:        map[string]string{"GOLINE": "bad"},
			getwd:      os.Getwd,
			wantCode:   1,
			wantStderr: "failed to parse GOLINE",
		},
		{
			name: "getwd error",
			env:  map[string]string{"GOLINE": "3"},
			getwd: func() (string, error) {
				return "", errors.New("no cwd")
			},
			wantCode:   1,
			wantStderr: "failed to get current directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			getwd := tt.getwd
			if getwd == nil {
				getwd = os.Getwd
			}

			exitCode := run(tt.args, testEnv(tt.env), getwd, &stdout, &stderr)
			if exitCode != tt.wantCode {
				t.Fatalf("run() exit = %d, want %d", exitCode, tt.wantCode)
			}
			if !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Fatalf("stderr = %q, want substring %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRunGenerationErrors(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		source     string
		env        map[string]string
		wantStderr string
	}{
		{
			name: "target file missing",
			env: map[string]string{
				"GOFILE":    "missing.go",
				"GOPACKAGE": "sample",
				"GOLINE":    "3",
			},
			wantStderr: "failed to create enum parser",
		},
		{
			name: "target line missing enum",
			source: `package sample

//go:generate enum_codegen
type Status int

const StatusReady Status = iota
`,
			env: map[string]string{
				"GOFILE":    "enum.go",
				"GOPACKAGE": "sample",
				"GOLINE":    "1",
			},
			wantStderr: "failed to get enum name",
		},
		{
			name: "no enum values",
			source: `package sample

//go:generate enum_codegen
type Status int

const Other = 1
`,
			env: map[string]string{
				"GOFILE":    "enum.go",
				"GOPACKAGE": "sample",
				"GOLINE":    "3",
			},
			wantStderr: "failed to parse enums",
		},
		{
			name: "forced undefined missing",
			args: []string{"-f"},
			source: `package sample

//go:generate enum_codegen
type Status int

const StatusReady Status = iota
`,
			env: map[string]string{
				"GOFILE":    "enum.go",
				"GOPACKAGE": "sample",
				"GOLINE":    "3",
			},
			wantStderr: "failed to validate enums",
		},
		{
			name: "generator error",
			source: `package sample

//go:generate enum_codegen
type Status int

const StatusReady Status = iota
`,
			env: map[string]string{
				"GOFILE":    "enum.go",
				"GOPACKAGE": "",
				"GOLINE":    "3",
			},
			wantStderr: "failed to generate code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.source != "" {
				writeCLITestFile(t, dir, "enum.go", tt.source)
			}

			var stdout, stderr bytes.Buffer
			exitCode := run(tt.args, testEnv(tt.env), func() (string, error) {
				return dir, nil
			}, &stdout, &stderr)

			if exitCode != 1 {
				t.Fatalf("run() exit = %d, want 1", exitCode)
			}
			if !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Fatalf("stderr = %q, want substring %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func testEnv(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}

func writeCLITestFile(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", name, err)
	}
}
