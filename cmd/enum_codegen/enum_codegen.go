package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/GeekchanskiY/enum_codegen/pkg/generator"
	"github.com/GeekchanskiY/enum_codegen/pkg/parser"
)

func main() {
	os.Exit(run(os.Args[1:], os.Getenv, os.Getwd, os.Stdout, os.Stderr))
}

func run(args []string, getenv func(string) string, getwd func() (string, error), stdout, stderr io.Writer) int {
	//
	// Flag parsing
	//

	flags := flag.NewFlagSet("enum_codegen", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(stdout, "Enum codegen by GeekchanskiY \n")
		_, _ = fmt.Fprintf(stdout, "Usage: enum_codegen [options]\n\n")
		_, _ = fmt.Fprintln(stdout, "Options:")

		flags.SetOutput(stdout)
		flags.PrintDefaults()
		flags.SetOutput(stderr)
	}

	help := flags.Bool("help", false, "show help")
	flags.BoolVar(help, "h", false, "show help")

	forceUndefined := flags.Bool("force-undefined", false, "force undefined enums")
	flags.BoolVar(forceUndefined, "f", false, "force undefined enums")

	if err := flags.Parse(args); err != nil {
		return 2
	}

	if *help {
		flags.Usage()
		return 0
	}

	//
	// ENV parsing
	//

	goFile := getenv("GOFILE")
	goPackage := getenv("GOPACKAGE")

	goLine, err := strconv.Atoi(getenv("GOLINE"))
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to parse GOLINE: %s\n", err)
		return 1
	}

	path, err := getwd()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to get current directory: %s\n", err)
		return 1
	}

	fullPath := filepath.Join(path, goFile)

	//
	// Enum parsing and generation
	//

	eParser, err := parser.New(path, fullPath, goLine)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to create enum parser: %s\n", err)
		return 1
	}

	enumName, err := eParser.GetEnumName()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to get enum name: %s\n", err)
		return 1
	}

	data, err := eParser.Parse()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to parse enums: %s\n", err)
		return 1
	}

	err = data.Validate(*forceUndefined)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to validate enums: %s\n", err)
		return 1
	}

	dataPath, err := generator.Generate(goFile, goPackage, path, enumName, data)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to generate code: %s\n", err)
		return 1
	}

	_, _ = fmt.Fprintf(stdout, "generated code to %s\n", dataPath)
	return 0
}
