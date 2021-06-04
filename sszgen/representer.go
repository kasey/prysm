package sszgen

import (
	"fmt"
	"go/ast"

	"github.com/prysmaticlabs/prysm/sszgen/types"
)

type Representer struct {
	index *PackageIndex
}

func NewRepresenter(pi *PackageIndex) *Representer {
	return &Representer{index: pi}
}

func (r *Representer) GetRepresentation(ts *TypeSpec) (types.ValRep, error) {
	typeSpec, err := r.index.GetStruct(ts.PackagePath, ts.Name)
	if err != nil {
		return nil, err
	}
	typeSpec.ValRep = &types.ValueContainer{}
	return expandRepresentation(typeSpec, r.index)
}

func expandRepresentation(ts *TypeSpec, index *PackageIndex) (types.ValRep, error) {
	switch ty := ts.TypeExpression.(type) {
	case *ast.StructType:
		vr := types.ValueContainer{
			Name: ts.Name,
			Contents: make([]types.ValRep, 0),
		}
		for _, f := range ty.Fields.List {
			// this filters out internal protobuf fields, but also serializers like us
			// can safely ignore unexported fields in general. We also ignore embedded
			// fields because I'm not sure if we should support them yet.
			if f.Names == nil || !ast.IsExported(f.Names[0].Name) {
				continue
			}
			s := &TypeSpec{
				Name: f.Names[0].Name,
				File: ts.File,
				PackageParser: ts.PackageParser,
				TypeExpression: f.Type,
				Tag: f.Tag.Value,
			}
			rep, err := expandRepresentation(s, index)
			if err != nil {
				return nil, err
			}
			vr.Contents = append(vr.Contents, rep)
		}
		return vr, nil
	case *ast.ArrayType:
		dims, err := extractSSZDimensions(ts.Tag)
		if err != nil {
			return nil, err
		}
		return expandArray(dims, ty, ts, index)
	case *ast.StarExpr:
		fmt.Printf("%v", ty)
	case *ast.SelectorExpr:
		fmt.Printf("%v", ty)
	case *ast.Ident:
		return expandPrimitive(ty, ts)
	}
	return nil, nil
}

func expandArray(dims []*SSZDimension, art *ast.ArrayType, ts *TypeSpec, index *PackageIndex) (types.ValRep, error) {
	if len(dims) == 0 {
		return nil, fmt.Errorf("Do not have dimension information for type %v", ts)
	}
	d := dims[0]
	var elv types.ValRep
	var err error
	switch elt := art.Elt.(type) {
	case *ast.ArrayType:
		elv, err = expandArray(dims[1:], elt, ts, index)
		if err != nil {
			return nil, err
		}
	default:
		elv, err = expandRepresentation(&TypeSpec{
			Name: ts.Name,
			File: ts.File,
			PackageParser: ts.PackageParser,
			TypeExpression: elt,
		}, index)
		if err != nil {
			return nil, err
		}
	}
	if d.IsVector() {
		return types.ValueVector{
			Name: ts.Name,
			ElementValue: elv,
			Size: d.VectorLen(),
		}, nil
	}
	if d.IsList() {
		return types.ValueList{
			Name: ts.Name,
			ElementValue: elv,
			MaxSize: d.ListLen(),
		}, nil
	}
	return nil, nil
}

func expandPrimitive(ident *ast.Ident, ts *TypeSpec) (types.ValRep, error) {
	switch ident.Name {
	case "bool":
		return types.ValueBool{Name: ts.Name}, nil
	case "byte":
		return types.ValueByte{Name: ts.Name}, nil
	case "uint8":
		return &types.ValueUint{Size: 8, Name: ts.Name}, nil
	case "uint16":
		return &types.ValueUint{Size: 16, Name: ts.Name}, nil
	case "uint32":
		return &types.ValueUint{Size: 32, Name: ts.Name}, nil
	case "uint64":
		return &types.ValueUint{Size: 64, Name: ts.Name}, nil
	case "uint128":
		return &types.ValueUint{Size: 128, Name: ts.Name}, nil
	case "uint256":
		return &types.ValueUint{Size: 256, Name: ts.Name}, nil
	}
	return nil, fmt.Errorf("Could not find primitive type matching %v", ident)
}
