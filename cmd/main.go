package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	// cfg := packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | packages.LoadSyntax}
	//pkgs, err := packages.Load(&cfg, "/Users/geekchanskiy/repos/enum_codegen/test_package")
	//if err != nil {
	//	panic(err)
	//}
	//
	//if packages.PrintErrors(pkgs) > 0 {
	//	panic(pkgs)
	//}
	//
	//for i, val := range pkgs[0].GoFiles {
	//	fmt.Println(i, val)
	//}
	//
	//fmt.Println(pkgs[0].Types.Scope().Lookup("Test").Type())

	fmt.Println(os.Args)

	fmt.Println(os.Getenv("GOARCH"))
	fmt.Println(os.Getenv("GOFILE"))
	fmt.Println(os.Getenv("GOOS"))
	fmt.Println(os.Getenv("GOLINE"))
	fmt.Println(os.Getenv("GOPACKAGE"))

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullPath := filepath.Join(path, os.Getenv("GOFILE"))

	fmt.Println(fullPath)

	fset := token.NewFileSet()
	data, err := parser.ParseFile(fset, fullPath, nil, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Scope.Lookup("Enum").Data)

	ast.Inspect(data, func(n ast.Node) bool {
		switch x := n.(type) {
		//case *ast.Ident:
		//	fmt.Println(x.Name)
		case *ast.TypeSpec:
			// TODO: get row number here somehow (?)
			fmt.Println(x.Name)
		}

		return true
	})
}
