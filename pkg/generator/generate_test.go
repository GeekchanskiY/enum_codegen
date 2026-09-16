package generator

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GeekchanskiY/enum_codegen/pkg/enum"
)

func TestGenerateCreatesCompilableCodeForMultipleEnums(t *testing.T) {
	dir := t.TempDir()
	writeGeneratorTestFile(t, dir, "go.mod", `module generated.test

go 1.23.0
`)
	writeGeneratorTestFile(t, dir, "model.go", `package sample

type Status int

const (
	StatusUnknown Status = iota
	StatusReady
)

type State int

const (
	Undefined State = iota
	StateActive
)
`)

	statusData := enum.Enum{
		{Name: "StatusUnknown", Value: 0, SnakeName: "unknown", Translate: "Unknown"},
		{Name: "StatusReady", Value: 1, SnakeName: `ready\value`, Translate: `Ready\translation`},
	}
	statusPath, err := Generate("status.kind.go", "sample", dir, "Status", statusData)
	if err != nil {
		t.Fatalf("Generate(Status) error = %v", err)
	}
	if !strings.HasSuffix(statusPath, "status.kind_Status__gen.go") {
		t.Fatalf("Generate(Status) path = %q, want status.kind_Status__gen.go suffix", statusPath)
	}

	stateData := enum.Enum{
		{Name: "Undefined", Value: 0, SnakeName: "undefined", Translate: "Undefined"},
		{Name: "StateActive", Value: 1, SnakeName: "state_active", Translate: "Active"},
	}
	if _, err = Generate("state.go", "sample", dir, "State", stateData); err != nil {
		t.Fatalf("Generate(State) error = %v", err)
	}

	writeGeneratorTestFile(t, dir, "generated_test.go", `package sample

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestGeneratedStatus(t *testing.T) {
	var status Status
	if err := status.Scan([]byte("ready\\value")); err != nil {
		t.Fatalf("Scan([]byte) error = %v", err)
	}
	if status != StatusReady {
		t.Fatalf("Scan([]byte) status = %v, want %v", status, StatusReady)
	}

	if err := status.Scan("missing"); err == nil {
		t.Fatal("Scan(missing) error = nil, want error without Undefined")
	}

	if StatusTags[StatusReady] != "ready\\value" {
		t.Fatalf("StatusTags[StatusReady] = %q", StatusTags[StatusReady])
	}
	if StatusTranslations[StatusReady] != "Ready\\translation" {
		t.Fatalf("StatusTranslations[StatusReady] = %q", StatusTranslations[StatusReady])
	}

	var _ driver.Valuer = status
	data, err := json.Marshal(StatusReady)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if string(data) != "\"ready\\\\value\"" {
		t.Fatalf("MarshalJSON() = %s", data)
	}
}

func TestGeneratedStateWithUndefined(t *testing.T) {
	var state State
	if err := state.Scan("missing"); err != nil {
		t.Fatalf("Scan(missing) error = %v", err)
	}
	if state != Undefined {
		t.Fatalf("Scan(missing) state = %v, want Undefined", state)
	}

	state = StateActive
	if err := state.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error = %v", err)
	}
	if state != Undefined {
		t.Fatalf("Scan(nil) state = %v, want Undefined", state)
	}

	var decoded State
	if err := json.Unmarshal([]byte("\"state_active\""), &decoded); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if decoded != StateActive {
		t.Fatalf("UnmarshalJSON() state = %v, want %v", decoded, StateActive)
	}

	if err := json.Unmarshal([]byte("\"missing\""), &decoded); err == nil {
		t.Fatal("UnmarshalJSON(missing) error = nil, want error")
	}
}
`)

	cmd := exec.Command("go", "test", ".")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated package go test failed: %v\n%s", err, output)
	}
}

func TestCompileTemplateReturnsWriterError(t *testing.T) {
	wantErr := errors.New("write failed")
	err := CompileTemplate(errorWriter{err: wantErr}, "sample", "Status", enum.Enum{
		{Name: "StatusReady", Value: 1, SnakeName: "ready", Translate: "Ready"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("CompileTemplate() error = %v, want %v", err, wantErr)
	}
}

func TestCompileTemplateReturnsParseError(t *testing.T) {
	err := compileTemplate(errorWriter{}, "{{", "sample", "Status", nil)
	if err == nil {
		t.Fatal("compileTemplate() error = nil, want parse error")
	}
}

func TestGenerateErrors(t *testing.T) {
	t.Run("template execution failure", func(t *testing.T) {
		_, err := Generate("status.go", "sample", t.TempDir(), "Status", enum.Enum{nil})
		if err == nil {
			t.Fatal("Generate() error = nil, want error")
		}
	})

	t.Run("invalid generated source", func(t *testing.T) {
		_, err := Generate("status.go", "sample", t.TempDir(), "Bad-Name", enum.Enum{
			{Name: "StatusReady", Value: 1, SnakeName: "ready", Translate: "Ready"},
		})
		if err == nil {
			t.Fatal("Generate() error = nil, want error")
		}
	})

	t.Run("write failure", func(t *testing.T) {
		missingDir := filepath.Join(t.TempDir(), "missing")
		_, err := Generate("status.go", "sample", missingDir, "Status", enum.Enum{
			{Name: "StatusReady", Value: 1, SnakeName: "ready", Translate: "Ready"},
		})
		if err == nil {
			t.Fatal("Generate() error = nil, want error")
		}
	})
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

func writeGeneratorTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", name, err)
	}

	return path
}
