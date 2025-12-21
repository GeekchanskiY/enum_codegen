package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/GeekchanskiY/enum_codegen/pkg/generator"
	"github.com/GeekchanskiY/enum_codegen/pkg/parser"
)

func main() {
	//
	// Flag parsing
	//

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stdout, "Enum codegen by GeekchanskiY \n")
		_, _ = fmt.Fprintf(os.Stdout, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintln(os.Stdout, "Options:")

		flag.PrintDefaults()
	}

	help := flag.Bool("help", false, "show help")
	flag.BoolVar(help, "h", false, "show help")

	forceUndefined := flag.Bool("force-undefined", false, "force undefined enums")
	flag.BoolVar(forceUndefined, "f", false, "force undefined enums")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	//
	// ENV parsing
	//

	goFile := os.Getenv("GOFILE")
	goPackage := os.Getenv("GOPACKAGE")

	goLine, err := strconv.Atoi(os.Getenv("GOLINE"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse GOLINE: %s\n", err)
		os.Exit(1)
	}

	path, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get current directory: %s\n", err)
		os.Exit(1)
	}

	fullPath := filepath.Join(path, goFile)

	//
	// Enum parsing and generation
	//

	eParser, err := parser.New(path, fullPath, goLine)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create enum parser: %s\n", err)
		os.Exit(1)
	}

	enumName, err := eParser.GetEnumName()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get enum name: %s\n", err)
		os.Exit(1)
	}

	data, err := eParser.Parse()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse enums: %s\n", err)
		os.Exit(1)
	}

	err = data.Validate(*forceUndefined)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to validate enums: %s\n", err)
		os.Exit(1)
	}

	dataPath, err := generator.Generate(goFile, goPackage, path, enumName, data)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to generate code: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("generated code to %s\n", dataPath)
}
