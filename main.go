package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const Template = `// Code generated via enum_codegen DO NOT EDIT.
package {{ .PackageName }}

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	_ sql.Scanner   = (*{{ .EnumName }})(nil)
	_ driver.Valuer = (*{{ .EnumName }})(nil)
	_ fmt.Stringer  = (*{{ .EnumName }})(nil)
	_ json.Marshaler = (*{{ .EnumName }})(nil)
)

var Tags = map[{{ .EnumName }}]string{
	{{- range .Enums }} 
		{{ .Name }}: "{{ .SnakeName }}",
	{{- end }}
}

var Types = map[string]{{ .EnumName }}{
	{{- range .Enums }} 
		"{{ .SnakeName }}": {{ .Name }},
	{{- end }}
}

var Translations = map[{{ .EnumName }}]string{
	{{- range .Enums }} 
		{{ .Name }}: "{{ .Translate }}",
	{{- end }}
}

func (t *{{ .EnumName }}) Scan(src any) error {
	value, ok := src.(string)
	
	if !ok {
		return errors.New("src is not string")
	}

	*t = Undefined
	if v, ok := Types[value]; ok {
		*t = v
	}

	return nil
}

func (t {{ .EnumName }}) Value() (driver.Value, error) {
	return Tags[t], nil
}

func (t {{ .EnumName }}) String() string {
	return Tags[t]
}

func (t {{ .EnumName }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *{{ .EnumName }}) UnmarshalJSON(data []byte) error {
	var (
		s string
		err error
	)

	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}

	if v, ok := Types[s]; ok {
		*t = v

		return nil
	}

	return fmt.Errorf("invalid status: %s", s)
}
`

type EnumData struct {
	Name      string
	Value     int64
	SnakeName string
	Translate string
}

func main() {
	targetLine, err := strconv.Atoi(os.Getenv("GOLINE"))
	if err != nil {
		panic(err)
	}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullPath := filepath.Join(path, os.Getenv("GOFILE"))

	fset := token.NewFileSet()

	data, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var enumName string

	ast.Inspect(data, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			start := fset.Position(n.Pos())

			// GOLINE is 1 line upper than n.Pos()
			if start.Line == targetLine+1 {
				fmt.Println("target found!", x.Name)
				enumName = x.Name.Name

				return false
			}
		}

		return true
	})

	if enumName == "" {
		panic("target not found")
	}

	// type parser
	conf := types.Config{}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	_, err = conf.Check(path, fset, []*ast.File{data}, info)
	if err != nil {
		panic(err)
	}

	neededDeclarations := make([]string, 0)

	// search for declared names of the type
	for ident, obj := range info.Defs {
		if obj != nil {
			// obj.Type().Underlying() - get base type
			objPathSplit := strings.Split(obj.Type().String(), ".")

			// base objects have different structure, filter them
			// may try obj.Type().Underlying() == obj.Type(), like I'm smart
			if len(objPathSplit) != 2 {
				continue
			}

			if objPathSplit[1] == enumName {
				neededDeclarations = append(neededDeclarations, ident.Name)
			}

		}
	}

	if len(neededDeclarations) == 0 {
		panic("enum declaration not found")
	}

	enums := make([]EnumData, 0, len(neededDeclarations))

	ast.Inspect(data, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			for _, spec := range x.Specs {
				if ts, ok := spec.(*ast.ValueSpec); ok {

					// TODO: research problem below
					// const EnumValue5, EnumValue6 Enum = 5, 6 will not be parsed for some reason
					if len(ts.Names) != 1 {
						continue
					}

					declaration := ts.Names[0].Obj
					if declaration == nil {
						continue
					}

					if slices.Contains(neededDeclarations, declaration.Name) {
						enumValue, err := strconv.ParseInt(fmt.Sprint(declaration.Data), 10, 64)
						if err != nil {
							panic(err)
						}

						translation := ""
						if ts.Doc != nil && len(ts.Doc.List) > 0 {
							translation = ts.Doc.List[0].Text
						}

						if translation = GetTranslationFromComment(translation); translation == "" {
							translation = CamelToSnake(declaration.Name)
						}

						enums = append(enums, EnumData{
							Name:      declaration.Name,
							Value:     enumValue,
							SnakeName: CamelToSnake(declaration.Name),
							Translate: translation,
						})
					}
				}
			}
		}

		return true
	})

	// force undefined value
	undefinedExists := false
	for _, v := range enums {
		if v.Name == "Undefined" {
			undefinedExists = true
			break
		}
	}
	if !undefinedExists {
		panic("must specify undefined value for enum")
	}

	// template generation
	tmpl, err := template.New("enum_code").Parse(Template)
	if err != nil {
		panic(err)
	}

	newFileName := strings.Split(os.Getenv("GOFILE"), ".")[0] + "_" + enumName + "__gen.go"

	dataPath := filepath.Join(path, newFileName)

	file, err := os.Create(dataPath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	err = tmpl.Execute(file, map[string]any{
		"PackageName": os.Getenv("GOPACKAGE"),
		"Enums":       enums,
		"EnumName":    enumName,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("generated code to %s\n", dataPath)
}

// CamelToSnake converts CamelCase or PascalCase with numbers to snake_case
func CamelToSnake(s string) string {
	re1 := regexp.MustCompile("([a-z0-9])([A-Z])")
	s = re1.ReplaceAllString(s, "${1}_${2}")

	re2 := regexp.MustCompile("([a-zA-Z])([0-9])")
	s = re2.ReplaceAllString(s, "${1}_${2}")

	re3 := regexp.MustCompile("([0-9])([a-zA-Z])")
	s = re3.ReplaceAllString(s, "${1}_${2}")

	return strings.ToLower(s)
}

// GetTranslationFromComment - get enum's translation from comment
func GetTranslationFromComment(comment string) string {
	re := regexp.MustCompile(`Translate="([^"]+)"`)
	match := re.FindStringSubmatch(comment)

	if len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}
