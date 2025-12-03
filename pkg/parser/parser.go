package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"

	"github.com/GeekchanskiY/enum_codegen/pkg/enum"
)

type EnumParser interface {
	GetEnumName() (string, error)
	Parse() (enum.Enum, error)
}

type enumParser struct {
	path, fullPath string
	goline         int

	// parser internals
	file    *ast.File
	fileset *token.FileSet
	config  types.Config
	info    *types.Info

	// computed values
	enumName string
}

var _ EnumParser = (*enumParser)(nil)

func New(path, fullPath string, goline int) (EnumParser, error) {
	file, fileset, err := getFile(fullPath)
	if err != nil {
		return nil, err
	}

	conf := types.Config{}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	_, err = conf.Check(path, fileset, []*ast.File{file}, info)
	if err != nil {
		return nil, err
	}

	return &enumParser{
		path:     path,
		fullPath: fullPath,
		goline:   goline,

		file:    file,
		fileset: fileset,
		config:  conf,
		info:    info,
	}, nil
}

func getFile(filename string) (*ast.File, *token.FileSet, error) {
	fileset := token.NewFileSet()

	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	return file, fileset, nil
}

func (p *enumParser) GetEnumName() (string, error) {
	if p.enumName != "" {
		return p.enumName, nil
	}

	ast.Inspect(p.file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			start := p.fileset.Position(n.Pos())

			// GOLINE is 1 line upper than n.Pos()
			if start.Line == p.goline+1 {
				p.enumName = x.Name.Name

				return false
			}
		}

		return true
	})

	if p.enumName == "" {
		return "", ErrTargetNotFound
	}

	return p.enumName, nil
}

func (p *enumParser) getVariableDeclarations() ([]string, error) {

	enumName, err := p.GetEnumName()
	if err != nil {
		return nil, err
	}

	declarations := make([]string, 0)

	// search for declared names of the type
	for ident, obj := range p.info.Defs {
		if obj != nil {
			// obj.Type().Underlying() - get base type
			objPathSplit := strings.Split(obj.Type().String(), ".")

			// base objects have different structure, filter them
			// may try obj.Type().Underlying() == obj.Type(), like I'm smart
			if len(objPathSplit) != 2 {
				continue
			}

			if ident.Name == enumName {
				continue
			}

			if objPathSplit[1] == enumName {
				declarations = append(declarations, ident.Name)
			}

		}
	}

	if len(declarations) == 0 {
		return nil, ErrTargetNotFound
	}

	return declarations, nil
}

func (p *enumParser) Parse() (enum.Enum, error) {
	neededDeclarations, err := p.getVariableDeclarations()
	if err != nil {
		return nil, err
	}

	enums := make([]*enum.Data, 0, len(neededDeclarations))

	ast.Inspect(p.file, func(n ast.Node) bool {
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
							return true
						}

						comment := ""
						if ts.Doc != nil && len(ts.Doc.List) > 0 {
							comment = ts.Doc.List[0].Text
						}

						translation := ""
						if translation = GetTranslationFromComment(comment); translation == "" {
							translation = CamelToSnake(declaration.Name)
						}

						stringName := ""
						if stringName = GetValueFromComment(comment); stringName == "" {
							stringName = CamelToSnake(declaration.Name)
						}

						enums = append(enums, &enum.Data{
							Name:      declaration.Name,
							Value:     enumValue,
							SnakeName: stringName,
							Translate: translation,
						})
					}
				}
			}
		}

		return true
	})

	if len(enums) == 0 {
		return nil, ErrTargetNotFound
	}

	if len(enums) != len(neededDeclarations) {
		for _, e := range neededDeclarations {
			if slices.IndexFunc(enums, func(data *enum.Data) bool { return data.Name == e }) == -1 {
				return nil, fmt.Errorf("%w: failed to get %s enum value", ErrParsingFailed, e)
			}
		}

		return nil, ErrParsingFailed
	}

	return enums, nil
}
