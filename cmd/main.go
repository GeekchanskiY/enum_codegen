package main

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"os"
)

func main() {
	cfg := packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | packages.LoadSyntax}
	pkgs, err := packages.Load(&cfg, "/Users/geekchanskiy/repos/enum_codegen/test_package")
	if err != nil {
		panic(err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		panic(pkgs)
	}

	for i, val := range pkgs[0].GoFiles {
		fmt.Println(i, val)
	}

	fmt.Println(pkgs[0].Types.Scope().Lookup("Test").Type())

	fmt.Println(os.Args)

	fmt.Println(os.Getenv("GOARCH"))
	fmt.Println(os.Getenv(""))
}
