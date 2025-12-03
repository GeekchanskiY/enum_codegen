package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GeekchanskiY/enum_codegen/pkg/enum"
)

func Generate(goFile, goPackage, path, enumName string, data enum.Enum) (string, error) {
	newFileName := strings.Split(goFile, ".")[0] + "_" + enumName + "__gen.go"

	dataPath := filepath.Join(path, newFileName)

	file, err := os.Create(dataPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create file: %s\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := file.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to close file: %s\n", err)
			os.Exit(1)
		}
	}()

	if err = CompileTemplate(file, goPackage, enumName, data); err != nil {
		return "", err
	}

	return dataPath, nil
}
