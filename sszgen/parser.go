package sszgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/prysmaticlabs/prysm/sszgen/types"
	"golang.org/x/tools/go/packages"
)

type TypeSpec struct {
	PackagePath    string
	Name           string
	TypeExpression ast.Expr
	File           *ast.File
	PackageParser  PackageParser
	ValRep         types.ValRep
	Tag string
}

type PackageParser interface {
	Imports() ([]*ast.ImportSpec, error)
	AllStructs() []*TypeSpec
	GetStruct(name string) (*TypeSpec, error)
}

type packageParser struct {
	packagePath string
	files map[string]*ast.File
}

func (pp *packageParser) Imports() ([]*ast.ImportSpec, error) {
	imports := make([]*ast.ImportSpec, 0)
	for _, f := range pp.files {
		for _, imp := range f.Imports {
			imports = append(imports, imp)
		}
	}
	return imports, nil
}

func (pp *packageParser) AllStructs() []*TypeSpec {
	structs := make([]*TypeSpec, 0)
	for _, f := range pp.files {
		for name, obj := range f.Scope.Objects {
			if obj.Kind != ast.Typ {
				continue
			}
			typeSpec, ok := obj.Decl.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			ts := &TypeSpec{
				Name:           name,
				TypeExpression: structType,
				File:           f,
				PackagePath:    pp.packagePath,
			}
			structs = append(structs, ts)
		}
	}
	return structs
}

func (pp *packageParser) GetStruct(name string) (*TypeSpec, error) {
	for _, f := range pp.files {
		for objName, obj := range f.Scope.Objects {
			if obj.Kind != ast.Typ {
				continue
			}
			typeSpec, ok := obj.Decl.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if name == objName {
				return &TypeSpec{
					Name:           objName,
					TypeExpression: structType,
					File:           f,
					PackageParser:  pp,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find struct named '%s' in package %s", name, pp.packagePath)
}

func NewPackageParser(packageName string) (*packageParser, error) {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, []string{packageName}...)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		if pkg.ID != packageName {
			continue
		}
		pp := &packageParser{packagePath: pkg.ID, files: make(map[string]*ast.File)}
		for _, f := range pkg.GoFiles {
			syn, err := parser.ParseFile(token.NewFileSet(), f, nil, parser.AllErrors)
			if err != nil {
				return nil, err
			}
			pp.files[f] = syn
		}
		return pp, nil
	}
	return nil, fmt.Errorf("Package named '%s' could not be loaded from the go build system. Please make sure the current folder contains the go.mod for the target package, or that its go.mod is in a parent directory", packageName)
}

