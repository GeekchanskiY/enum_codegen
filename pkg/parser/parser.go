package parser

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
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
	file         *ast.File
	packageFiles []*ast.File
	fileset      *token.FileSet
	info         *types.Info

	// computed values
	enumName string
	enumType types.Type
}

var _ EnumParser = (*enumParser)(nil)

func New(path, fullPath string, goline int) (EnumParser, error) {
	file, packageFiles, fileset, err := getPackageFiles(path, fullPath)
	if err != nil {
		return nil, err
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	_, err = conf.Check(path, fileset, packageFiles, info)
	if err != nil {
		return nil, err
	}

	return &enumParser{
		path:         path,
		fullPath:     fullPath,
		goline:       goline,
		packageFiles: packageFiles,

		file:    file,
		fileset: fileset,
		info:    info,
	}, nil
}

func getPackageFiles(path, fullPath string) (*ast.File, []*ast.File, *token.FileSet, error) {
	fileset := token.NewFileSet()
	packageFiles := make([]*ast.File, 0)

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, nil, err
	}

	var targetFile *ast.File
	targetPath := filepath.Clean(fullPath)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}

		filePath := filepath.Join(path, name)
		file, err := parser.ParseFile(fileset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, nil, nil, err
		}

		if filepath.Clean(filePath) == targetPath {
			targetFile = file
		}

		packageFiles = append(packageFiles, file)
	}

	if targetFile == nil {
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrTargetNotFound, fullPath)
	}

	return targetFile, packageFiles, fileset, nil
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
				if obj, ok := p.info.Defs[x.Name].(*types.TypeName); ok {
					p.enumType = obj.Type()
				}

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

func (p *enumParser) Parse() (enum.Enum, error) {
	enumName, err := p.GetEnumName()
	if err != nil {
		return nil, err
	}
	if p.enumType == nil {
		return nil, fmt.Errorf("%w: %s enum type", ErrTargetNotFound, enumName)
	}

	enums := make([]*enum.Data, 0)

	for _, file := range p.packageFiles {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				comment := getCommentText(valueSpec.Doc, valueSpec.Comment)
				for _, name := range valueSpec.Names {
					obj, ok := p.info.Defs[name].(*types.Const)
					if !ok || !types.Identical(obj.Type(), p.enumType) {
						continue
					}

					if obj.Val().Kind() != constant.Int {
						return nil, fmt.Errorf("%w: failed to get %s enum value", ErrParsingFailed, name.Name)
					}

					enumValue, ok := constant.Int64Val(obj.Val())
					if !ok {
						return nil, fmt.Errorf("%w: failed to get %s enum value", ErrParsingFailed, name.Name)
					}

					translation := GetTranslationFromComment(comment)
					if translation == "" {
						translation = CamelToSnake(name.Name)
					}

					stringName := GetValueFromComment(comment)
					if stringName == "" {
						stringName = CamelToSnake(name.Name)
					}

					enums = append(enums, &enum.Data{
						Name:      name.Name,
						Value:     enumValue,
						SnakeName: stringName,
						Translate: translation,
					})
				}
			}
		}
	}

	if len(enums) == 0 {
		return nil, ErrTargetNotFound
	}

	return enums, nil
}

func getCommentText(groups ...*ast.CommentGroup) string {
	comment := ""
	for _, group := range groups {
		if group != nil {
			comment += group.Text()
		}
	}

	return comment
}
