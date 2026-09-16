package generator

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/GeekchanskiY/enum_codegen/pkg/enum"
)

func Generate(goFile, goPackage, path, enumName string, data enum.Enum) (string, error) {
	newFileName := strings.TrimSuffix(goFile, filepath.Ext(goFile)) + "_" + enumName + "__gen.go"

	dataPath := filepath.Join(path, newFileName)

	var raw bytes.Buffer
	if err := CompileTemplate(&raw, goPackage, enumName, data); err != nil {
		return "", err
	}

	formatted, err := format.Source(raw.Bytes())
	if err != nil {
		return "", err
	}

	if err = os.WriteFile(dataPath, formatted, 0o644); err != nil {
		return "", err
	}

	return dataPath, nil
}
